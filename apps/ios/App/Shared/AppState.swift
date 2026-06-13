import BlocklistStore
import FilterEngine
import Foundation
import SharedModels

@MainActor
@Observable
final class AppState {
    var filterConfig = FilterConfig()
    var stats = FilterStats()
    var samples: [SampleEntry] = []
    var extensionEnabled = false
    var onboardingCompleted = false
    var userMode: UserMode = .standard
    var lastTraceID = ""
    var isPro: Bool { EntitlementManager.shared.isPro }

    private let store = BlocklistStore()
    private var engine = HybridFilterEngine()

    init() {
        loadUserDefaults()
    }

    func loadAll() async {
        do {
            filterConfig = try await store.loadConfig()
            stats = try await store.loadStats()
            if let modelData = try await store.loadBayesModel() {
                engine.loadBayesModel(from: modelData)
            }
            if let url = try? await store.coreMLCompiledURL(),
               FileManager.default.fileExists(atPath: url.path) {
                engine.loadCoreML(from: url)
            }
        } catch {
            ErrorPresenter.shared.present(MGError.store(.containerUnavailable))
        }
    }

    func saveConfig() async {
        do {
            try await store.saveConfig(filterConfig)
            if let modelData = engine.exportBayesModel() {
                try await store.saveBayesModel(modelData)
            }
            AnalyticsManager.shared.track(.settingsChanged(key: "config", value: "saved"))
        } catch {
            ErrorPresenter.shared.present(MGError.filter(.ruleSaveFailed))
        }
    }

    func submitSample(text: String, label: MessageCategory) async {
        engine.trainSample(text: text, category: label)
        samples.insert(SampleEntry(text: text, label: label), at: 0)
        if let modelData = engine.exportBayesModel() {
            try? await store.saveBayesModel(modelData)
        }
        AnalyticsManager.shared.track(.sampleSubmitted(label: label.rawValue))
    }

    func markExtensionEnabled(_ enabled: Bool) {
        extensionEnabled = enabled
        UserDefaults.standard.set(enabled, forKey: AppConstants.UserDefaultsKeys.extensionEnabled)
        AnalyticsManager.shared.track(.extensionEnabledChanged(enabled: enabled))
    }

    func completeOnboarding() {
        onboardingCompleted = true
        UserDefaults.standard.set(true, forKey: AppConstants.UserDefaultsKeys.onboardingCompleted)
        AnalyticsManager.shared.track(.onboardingCompleted)
    }

    func refreshStats() async {
        stats = (try? await store.loadStats()) ?? stats
    }

    private func loadUserDefaults() {
        let defaults = UserDefaults.standard
        extensionEnabled = defaults.bool(forKey: AppConstants.UserDefaultsKeys.extensionEnabled)
        onboardingCompleted = defaults.bool(forKey: AppConstants.UserDefaultsKeys.onboardingCompleted)
        if let mode = defaults.string(forKey: AppConstants.UserDefaultsKeys.userMode),
           let parsed = UserMode(rawValue: mode) {
            userMode = parsed
        }
    }
}
