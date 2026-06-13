import BlocklistStore
import Foundation
import SharedModels

/// Syncs rules and model metadata from backend to App Group.
actor SyncService {
    private let client = APIClient.shared
    private let store = BlocklistStore()

    func syncRules() async throws {
        struct RulesMeta: Decodable {
            let version: String
            let checksum: String?
            let keywords_block: [String]?
        }
        var config = try await store.loadConfig()
        let etag = config.rulesChecksum.isEmpty ? nil : config.rulesChecksum
        let result = try await client.fetch(
            APIEndpoint(path: "/api/v1/rules/latest"),
            ifNoneMatch: etag
        )
        switch result {
        case .notModified:
            MGLogger.sync.info("Rules unchanged \(config.rulesVersion)")
            return
        case .data(let data):
            let meta = try JSONDecoder().decode(RulesMeta.self, from: data)
            if meta.version == config.rulesVersion,
               let checksum = meta.checksum, checksum == config.rulesChecksum {
                return
            }
            if let keywords = meta.keywords_block {
                let remoteRules = keywords.map { BlockRule(type: .keywordBlock, pattern: $0, priority: 0) }
                let custom = config.rules.filter { $0.type != .keywordBlock || $0.priority > 0 }
                config.rules = custom + remoteRules
            }
            config.rulesVersion = meta.version
            config.rulesChecksum = meta.checksum ?? meta.version
            try await store.saveConfig(config)
            MGLogger.sync.info("Rules synced version \(meta.version)")
        }
    }

    func reportHealth(latencyMs: Double, layer: String) async {
        MGLogger.sync.info("health layer=\(layer) latency=\(latencyMs)ms")
    }
}
