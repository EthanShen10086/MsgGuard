import DesignSystem
import SharedModels
import SwiftUI

struct SettingsView: View {
    @Environment(AppState.self) private var appState

    var body: some View {
        NavigationStack {
            Form {
                Section(String(localized: "settings.extension")) {
                    Toggle(String(localized: "settings.extensionEnabled"), isOn: Binding(
                        get: { appState.extensionEnabled },
                        set: { appState.markExtensionEnabled($0) }
                    ))
                }
                Section(String(localized: "settings.privacy")) {
                    Toggle(String(localized: "settings.cloudLLM"), isOn: Binding(
                        get: { appState.filterConfig.cloudLLMEnabled },
                        set: {
                            guard EntitlementManager.shared.hasEntitlement(.cloudLLM) else { return }
                            appState.filterConfig.cloudLLMEnabled = $0
                            Task { await appState.saveConfig() }
                        }
                    ))
                    .disabled(!EntitlementManager.shared.hasEntitlement(.cloudLLM))
                    if !EntitlementManager.shared.hasEntitlement(.cloudLLM) {
                        Text(String(localized: "settings.cloudLLMProHint"))
                            .font(.caption)
                            .foregroundStyle(.secondary)
                    }
                }
                Section(String(localized: "settings.sync")) {
                    Toggle(String(localized: "settings.iCloudSync"), isOn: Binding(
                        get: { appState.filterConfig.iCloudSyncEnabled },
                        set: {
                            guard EntitlementManager.shared.hasEntitlement(.autoSync) else { return }
                            appState.filterConfig.iCloudSyncEnabled = $0
                            Task { await appState.saveConfig() }
                        }
                    ))
                    .disabled(!EntitlementManager.shared.hasEntitlement(.autoSync))
                    if !EntitlementManager.shared.hasEntitlement(.autoSync) {
                        Text(String(localized: "settings.iCloudSyncProHint"))
                            .font(.caption)
                            .foregroundStyle(.secondary)
                    }
                }
                Section(String(localized: "settings.callDirectory")) {
                    Text(String(localized: "settings.callDirectoryHint"))
                        .font(.caption)
                        .foregroundStyle(.secondary)
                }
                Section(String(localized: "settings.accessibility")) {
                    Picker(String(localized: "settings.userMode"), selection: Binding(
                        get: { appState.userMode },
                        set: {
                            appState.userMode = $0
                            UserDefaults.standard.set($0.rawValue, forKey: AppConstants.UserDefaultsKeys.userMode)
                        }
                    )) {
                        Text("Standard").tag(UserMode.standard)
                        Text("Elder").tag(UserMode.elder)
                    }
                }
                Section(String(localized: "settings.subscription")) {
                    if appState.isPro {
                        Text(String(localized: "settings.proActive"))
                        Button(String(localized: "subscription.manage")) {
                            Task { await StoreManager.shared.openManageSubscriptions() }
                        }
                    } else {
                        NavigationLink(String(localized: "settings.upgradePro")) {
                            SubscriptionView()
                        }
                    }
                    Button(String(localized: "subscription.restore")) {
                        Task { await StoreManager.shared.restorePurchases() }
                    }
                }
                Section {
                    Link(String(localized: "settings.support"), destination: URL(string: "https://msgguard.app/support")!)
                    Link(String(localized: "settings.privacyPolicy"), destination: URL(string: "https://msgguard.app/privacy")!)
                }
            }
            .navigationTitle(String(localized: "tab.settings"))
        }
    }
}
