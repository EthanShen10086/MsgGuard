import StoreKit
import SwiftUI

struct SubscriptionView: View {
    @State private var store = StoreManager.shared

    var body: some View {
        List {
            Section(String(localized: "subscription.features")) {
                Text(String(localized: "subscription.feature.cloud"))
                Text(String(localized: "subscription.feature.rules"))
                Text(String(localized: "subscription.feature.stats"))
            }
            Section {
                if store.isLoading {
                    ProgressView()
                }
                ForEach(store.products, id: \.id) { product in
                    Button("\(product.displayName) — \(product.displayPrice)") {
                        Task { await purchase(product) }
                    }
                }
            }
            Section {
                Button(String(localized: "subscription.restore")) {
                    Task { await store.restorePurchases() }
                }
            }
        }
        .navigationTitle(String(localized: "settings.upgradePro"))
        .task { await store.loadProducts() }
    }

    private func purchase(_ product: Product) async {
        do {
            _ = try await store.purchase(product)
        } catch {
            ErrorPresenter.shared.present(MGError.generic(error.localizedDescription))
        }
    }
}
