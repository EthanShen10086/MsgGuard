import BlocklistStore
import SharedModels
import WidgetKit

/// Downloads model artifacts from backend and writes to App Group.
actor ModelUpdateService {
    private let client = APIClient.shared
    private let store = BlocklistStore()

    struct ModelMeta: Decodable {
        let version: String
        let checksum: String
        let url: String
    }

    func checkAndUpdate() async throws {
        let meta: ModelMeta = try await client.request(
            APIEndpoint(path: "/api/v1/models/latest?locale=zh-Hans")
        )
        let config = try await store.loadConfig()
        if config.modelVersion == meta.version {
            MGLogger.sync.info("Model up to date \(meta.version)")
            return
        }
        MGLogger.sync.info("Updating model to \(meta.version)")
        // Store version; actual .mlmodel download when bundled in CDN
        var updated = config
        updated.modelVersion = meta.version
        try await store.saveConfig(updated)
        WidgetCenter.shared.reloadAllTimelines()
    }
}
