import Foundation
import SharedModels

/// Syncs filter config via iCloud Key-Value storage (Pro autoSync entitlement).
@MainActor
final class CloudSyncService {
    static let shared = CloudSyncService()

    private let store = NSUbiquitousKeyValueStore.default
    private let configKey = "filter_config"

    private init() {
        NotificationCenter.default.addObserver(
            self,
            selector: #selector(cloudChanged),
            name: NSUbiquitousKeyValueStore.didChangeExternallyNotification,
            object: store
        )
        store.synchronize()
    }

    func pushConfig(_ config: FilterConfig) {
        guard let data = try? JSONEncoder().encode(config) else { return }
        store.set(data, forKey: configKey)
        store.synchronize()
        MGLogger.sync.info("iCloud config pushed")
    }

    func pullConfig() -> FilterConfig? {
        guard let data = store.data(forKey: configKey) else { return nil }
        return try? JSONDecoder().decode(FilterConfig.self, from: data)
    }

    @objc private func cloudChanged(_ note: Notification) {
        MGLogger.sync.info("iCloud external change received")
    }
}
