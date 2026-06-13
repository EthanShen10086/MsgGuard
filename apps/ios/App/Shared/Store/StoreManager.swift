import Foundation
import OSLog
import StoreKit
import UIKit

@MainActor
@Observable
final class StoreManager {
    static let shared = StoreManager()

    private(set) var products: [Product] = []
    private(set) var purchasedProductIDs: Set<String> = []
    private(set) var isLoading = false

    private let logger = Logger(subsystem: "com.ethanshen.msgguard", category: "store")
    private var updateListenerTask: Task<Void, Never>?

    static let productIDs = [
        "com.ethanshen.msgguard.pro.monthly",
        "com.ethanshen.msgguard.pro.yearly",
    ]

    static let subscriptionGroupID = "MG000001-0001"

    var isPro: Bool {
        !self.purchasedProductIDs.isEmpty
    }

    private init() {
        self.updateListenerTask = self.listenForTransactions()
        Task { await self.refreshPurchaseState() }
    }

    func loadProducts() async {
        self.isLoading = true
        defer { isLoading = false }
        do {
            self.products = try await Product.products(for: Self.productIDs)
                .sorted { $0.price < $1.price }
            self.logger.info("Loaded \(self.products.count) products")
        } catch {
            self.logger.error("load products failed: \(error.localizedDescription)")
            AnalyticsManager.shared.track(.error(domain: "store", code: "load_products_failed"))
        }
    }

    func purchase(_ product: Product) async throws -> Bool {
        let result = try await product.purchase()
        switch result {
        case let .success(verification):
            let transaction = try checkVerified(verification)
            self.purchasedProductIDs.insert(transaction.productID)
            await transaction.finish()
            EntitlementManager.shared.syncFromStoreKit(purchasedProductIDs: self.purchasedProductIDs)
            AnalyticsManager.shared.track(.purchaseCompleted(productId: product.id))
            return true
        case .userCancelled:
            AnalyticsManager.shared.track(.purchaseCancelled(productId: product.id))
            return false
        case .pending:
            return false
        @unknown default:
            return false
        }
    }

    func restorePurchases() async {
        var restored = 0
        for await result in Transaction.currentEntitlements {
            if let transaction = try? checkVerified(result) {
                self.purchasedProductIDs.insert(transaction.productID)
                restored += 1
            }
        }
        await self.refreshSubscriptionStatus()
        EntitlementManager.shared.syncFromStoreKit(purchasedProductIDs: self.purchasedProductIDs)
        AnalyticsManager.shared.track(.restoreCompleted(count: restored))
    }

    func refreshPurchaseState() async {
        for await result in Transaction.currentEntitlements {
            if let transaction = try? checkVerified(result) {
                self.purchasedProductIDs.insert(transaction.productID)
            }
        }
        await self.refreshSubscriptionStatus()
        EntitlementManager.shared.syncFromStoreKit(purchasedProductIDs: self.purchasedProductIDs)
    }

    func refreshSubscriptionStatus() async {
        guard let statuses = try? await Product.SubscriptionInfo.status(for: Self.subscriptionGroupID) else {
            return
        }
        for status in statuses {
            guard let transaction = try? checkVerified(status.transaction) else { continue }
            switch status.state {
            case .subscribed, .inGracePeriod:
                self.purchasedProductIDs.insert(transaction.productID)
            case .expired, .revoked:
                self.purchasedProductIDs.remove(transaction.productID)
            default:
                break
            }
        }
        EntitlementManager.shared.syncFromStoreKit(purchasedProductIDs: self.purchasedProductIDs)
    }

    func openManageSubscriptions() async {
        guard let scene = UIApplication.shared.connectedScenes
            .compactMap({ $0 as? UIWindowScene })
            .first
        else { return }
        try? await AppStore.showManageSubscriptions(in: scene)
    }

    private func listenForTransactions() -> Task<Void, Never> {
        Task.detached {
            for await result in Transaction.updates {
                guard case let .verified(transaction) = result else { continue }
                await MainActor.run {
                    self.purchasedProductIDs.insert(transaction.productID)
                    EntitlementManager.shared.syncFromStoreKit(purchasedProductIDs: self.purchasedProductIDs)
                }
                await transaction.finish()
                await self.refreshSubscriptionStatus()
            }
        }
    }

    private func checkVerified<T>(_ result: VerificationResult<T>) throws -> T {
        switch result {
        case .unverified:
            throw StoreError.verificationFailed
        case let .verified(safe):
            return safe
        }
    }
}

enum StoreError: Error, LocalizedError {
    case verificationFailed

    var errorDescription: String? {
        switch self {
        case .verificationFailed: "Transaction verification failed."
        }
    }
}
