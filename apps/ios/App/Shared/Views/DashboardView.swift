import DesignSystem
import SwiftUI

struct DashboardView: View {
    @Environment(AppState.self) private var appState

    var body: some View {
        NavigationStack {
            ScrollView {
                VStack(spacing: 16) {
                    MGCard {
                        VStack(alignment: .leading, spacing: 8) {
                            Label(
                                appState.extensionEnabled ? String(localized: "dashboard.protected") : String(localized: "dashboard.notProtected"),
                                systemImage: appState.extensionEnabled ? "checkmark.shield.fill" : "exclamationmark.shield.fill"
                            )
                            .font(.headline)
                            .foregroundStyle(appState.extensionEnabled ? .green : .orange)
                            Text(String(localized: "dashboard.blockedToday \(appState.stats.blockedToday)"))
                            Text(String(localized: "dashboard.blockedTotal \(appState.stats.blockedTotal)"))
                                .foregroundStyle(.secondary)
                        }
                        .frame(maxWidth: .infinity, alignment: .leading)
                    }
                    MGPrimaryButton(String(localized: "dashboard.refreshStats")) {
                        Task { await appState.refreshStats() }
                    }
                    MGPrimaryButton(String(localized: "dashboard.manageExtension")) {
                        if let url = URL(string: UIApplication.openSettingsURLString) {
                            UIApplication.shared.open(url)
                        }
                    }
                }
                .padding()
            }
            .navigationTitle(String(localized: "tab.home"))
            .task { await appState.refreshStats() }
        }
    }
}
