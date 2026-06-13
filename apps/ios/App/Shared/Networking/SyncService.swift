import BlocklistStore
import SharedModels

/// Syncs rules and model metadata from backend to App Group.
actor SyncService {
    private let client = APIClient.shared
    private let store = BlocklistStore()

    func syncRules() async throws {
        struct RulesMeta: Decodable {
            let version: String
            let keywords_block: [String]?
        }
        let meta: RulesMeta = try await client.request(APIEndpoint(path: "/api/v1/rules/latest"))
        var config = try await store.loadConfig()
        if let keywords = meta.keywords_block {
            let rules = keywords.map { BlockRule(type: .keywordBlock, pattern: $0, priority: 0) }
            config.rules.append(contentsOf: rules)
        }
        try await store.saveConfig(config)
        MGLogger.sync.info("Rules synced version \(meta.version)")
    }

    func reportHealth(latencyMs: Double, layer: String) async {
        MGLogger.sync.info("health layer=\(layer) latency=\(latencyMs)ms")
    }
}
