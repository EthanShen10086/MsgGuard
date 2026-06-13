import Foundation
import SharedModels

public actor BlocklistStore {
    private let appGroupID: String
    private let encoder = JSONEncoder()
    private let decoder = JSONDecoder()

    public init(appGroupID: String = AppConstants.appGroupID) {
        self.appGroupID = appGroupID
        encoder.outputFormatting = [.prettyPrinted, .sortedKeys]
    }

    public func loadConfig() async throws -> FilterConfig {
        let url = try configURL()
        guard FileManager.default.fileExists(atPath: url.path) else {
            return FilterConfig()
        }
        let data = try Data(contentsOf: url)
        return try decoder.decode(FilterConfig.self, from: data)
    }

    public func saveConfig(_ config: FilterConfig) async throws {
        let url = try configURL()
        try FileManager.default.createDirectory(at: url.deletingLastPathComponent(), withIntermediateDirectories: true)
        let data = try encoder.encode(config)
        try data.write(to: url, options: .atomic)
    }

    public func loadStats() async throws -> FilterStats {
        let url = try statsURL()
        guard FileManager.default.fileExists(atPath: url.path) else {
            return FilterStats()
        }
        let data = try Data(contentsOf: url)
        return try decoder.decode(FilterStats.self, from: data)
    }

    public func saveStats(_ stats: FilterStats) async throws {
        let url = try statsURL()
        try FileManager.default.createDirectory(at: url.deletingLastPathComponent(), withIntermediateDirectories: true)
        let data = try encoder.encode(stats)
        try data.write(to: url, options: .atomic)
    }

    public func loadBayesModel() async throws -> Data? {
        let url = try bayesURL()
        guard FileManager.default.fileExists(atPath: url.path) else { return nil }
        return try Data(contentsOf: url)
    }

    public func saveBayesModel(_ data: Data) async throws {
        let url = try bayesURL()
        try FileManager.default.createDirectory(at: url.deletingLastPathComponent(), withIntermediateDirectories: true)
        try data.write(to: url, options: .atomic)
    }

    public func incrementBlocked(category: MessageCategory) async throws {
        var stats = try await loadStats()
        let calendar = Calendar.current
        if !calendar.isDateInToday(stats.lastUpdated) {
            stats.blockedToday = 0
        }
        stats.blockedToday += 1
        stats.blockedTotal += 1
        stats.byCategory[category.rawValue, default: 0] += 1
        stats.lastUpdated = Date()
        try await saveStats(stats)
        UserDefaults(suiteName: appGroupID)?.set(stats.blockedToday, forKey: AppConstants.UserDefaultsKeys.blockedCountToday)
        UserDefaults(suiteName: appGroupID)?.set(stats.blockedTotal, forKey: AppConstants.UserDefaultsKeys.blockedCountTotal)
    }

    private func containerURL() throws -> URL {
        guard let url = FileManager.default.containerURL(forSecurityApplicationGroupIdentifier: appGroupID) else {
            throw BlocklistStoreError.containerUnavailable
        }
        return url
    }

    private func configURL() throws -> URL {
        try containerURL().appendingPathComponent(AppConstants.AppGroupFiles.rulesManifest)
    }

    private func statsURL() throws -> URL {
        try containerURL().appendingPathComponent(AppConstants.AppGroupFiles.statsSnapshot)
    }

    private func bayesURL() throws -> URL {
        try containerURL().appendingPathComponent(AppConstants.AppGroupFiles.bayesModel)
    }
}

public enum BlocklistStoreError: Error {
    case containerUnavailable
}
