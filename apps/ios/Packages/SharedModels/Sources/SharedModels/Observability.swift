import Foundation
import OSLog

public enum MGLogger {
    public static let filter = Logger(subsystem: AppConstants.bundleID, category: "filter")
    public static let sync = Logger(subsystem: AppConstants.bundleID, category: "sync")
    public static let network = Logger(subsystem: AppConstants.bundleID, category: "network")
    public static let feedback = Logger(subsystem: AppConstants.bundleID, category: "feedback")
    public static let app = Logger(subsystem: AppConstants.bundleID, category: "app")
}

public enum PerformanceMonitor {
    private static let log = OSLog(subsystem: AppConstants.bundleID, category: "performance")

    public static func beginFilter() -> OSSignpostID {
        let id = OSSignpostID(log: log)
        os_signpost(.begin, log: log, name: "FilterMessage", signpostID: id)
        return id
    }

    public static func endFilter(_ id: OSSignpostID, latencyMs: Double) {
        os_signpost(.end, log: log, name: "FilterMessage", signpostID: id, "latency=%.2fms", latencyMs)
        recordAggregate(latencyMs: latencyMs)
    }

    private static func recordAggregate(latencyMs: Double) {
        guard let url = FileManager.default.containerURL(
            forSecurityApplicationGroupIdentifier: AppConstants.appGroupID
        )?.appendingPathComponent("filter_perf.json") else { return }
        var stats: [String: Double] = [:]
        if let data = try? Data(contentsOf: url),
           let decoded = try? JSONDecoder().decode([String: Double].self, from: data) {
            stats = decoded
        }
        let count = (stats["count"] ?? 0) + 1
        let total = (stats["total_ms"] ?? 0) + latencyMs
        let max = max(stats["max_ms"] ?? 0, latencyMs)
        stats = ["count": count, "total_ms": total, "max_ms": max, "mean_ms": total / count]
        if let encoded = try? JSONEncoder().encode(stats) {
            try? encoded.write(to: url, options: .atomic)
        }
    }

    public static func loadAggregateStats() -> [String: Double] {
        guard let url = FileManager.default.containerURL(
            forSecurityApplicationGroupIdentifier: AppConstants.appGroupID
        )?.appendingPathComponent("filter_perf.json"),
              let data = try? Data(contentsOf: url),
              let stats = try? JSONDecoder().decode([String: Double].self, from: data) else {
            return [:]
        }
        return stats
    }
}
