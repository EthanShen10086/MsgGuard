import Foundation
import SharedModels

public enum SyncConfigLoader {
    public static func loadConfig(appGroupID: String = AppConstants.appGroupID) -> FilterConfig {
        guard let url = containerURL(appGroupID)?.appendingPathComponent(AppConstants.AppGroupFiles.rulesManifest),
              let data = try? Data(contentsOf: url),
              let config = try? JSONDecoder().decode(FilterConfig.self, from: data) else {
            return FilterConfig()
        }
        return config
    }

    public static func loadBayesModel(appGroupID: String = AppConstants.appGroupID) -> Data? {
        guard let url = containerURL(appGroupID)?.appendingPathComponent(AppConstants.AppGroupFiles.bayesModel) else {
            return nil
        }
        return try? Data(contentsOf: url)
    }

    public static func coreMLModelURL(appGroupID: String = AppConstants.appGroupID) -> URL? {
        containerURL(appGroupID)?.appendingPathComponent(AppConstants.AppGroupFiles.coreMLModel)
    }

    public static func saveStatsIncrement(category: MessageCategory, appGroupID: String = AppConstants.appGroupID) {
        guard let url = containerURL(appGroupID)?.appendingPathComponent(AppConstants.AppGroupFiles.statsSnapshot) else {
            return
        }
        var stats = FilterStats()
        if let data = try? Data(contentsOf: url),
           let decoded = try? JSONDecoder().decode(FilterStats.self, from: data) {
            stats = decoded
        }
        if !Calendar.current.isDateInToday(stats.lastUpdated) {
            stats.blockedToday = 0
        }
        stats.blockedToday += 1
        stats.blockedTotal += 1
        stats.byCategory[category.rawValue, default: 0] += 1
        stats.lastUpdated = Date()
        if let data = try? JSONEncoder().encode(stats) {
            try? data.write(to: url, options: .atomic)
        }
        UserDefaults(suiteName: appGroupID)?.set(stats.blockedToday, forKey: AppConstants.UserDefaultsKeys.blockedCountToday)
    }

    private static func containerURL(_ appGroupID: String) -> URL? {
        FileManager.default.containerURL(forSecurityApplicationGroupIdentifier: appGroupID)
    }
}
