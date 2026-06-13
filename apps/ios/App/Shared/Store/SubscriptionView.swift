import StoreKit
import SwiftUI

struct SubscriptionView: View {
    @State private var products: [Product] = []
    @State private var isLoading = false

    var body: some View {
        List {
            Section(String(localized: "subscription.features")) {
                Text(String(localized: "subscription.feature.cloud"))
                Text(String(localized: "subscription.feature.rules"))
                Text(String(localized: "subscription.feature.stats"))
            }
            Section {
                if isLoading {
                    ProgressView()
                }
                ForEach(products, id: \.id) { product in
                    Button("\(product.displayName) — \(product.displayPrice)") {
                        Task { await purchase(product) }
                    }
                }
            }
        }
        .navigationTitle(String(localized: "settings.upgradePro"))
        .task { await loadProducts() }
    }

    private func loadProducts() async {
        isLoading = true
        defer { isLoading = false }
        do {
            products = try await Product.products(for: ["com.ethanshen.msgguard.pro.monthly", "com.ethanshen.msgguard.pro.yearly"])
        } catch {
            ErrorPresenter.shared.present(MGError.generic(error.localizedDescription))
        }
    }

    private func purchase(_ product: Product) async {
        do {
            let result = try await product.purchase()
            if case let .success(verification) = result,
               case .verified = verification {
                EntitlementManager.shared.grantPro(source: .appStoreIAP, transactionId: product.id)
                AnalyticsManager.shared.track(.purchaseCompleted(productId: product.id))
            }
        } catch {
            ErrorPresenter.shared.present(MGError.generic(error.localizedDescription))
        }
    }
}
