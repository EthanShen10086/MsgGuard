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
    }
}
