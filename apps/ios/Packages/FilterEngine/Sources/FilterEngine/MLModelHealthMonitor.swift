import Foundation

public struct MLModelHealthMonitor: Sendable {
    private var totalInferences: Int = 0
    private var totalLatencyMs: Double = 0

    public init() {}

    public mutating func record(latencyMs: Double) {
        totalInferences += 1
        totalLatencyMs += latencyMs
    }

    public var meanLatencyMs: Double {
        guard totalInferences > 0 else { return 0 }
        return totalLatencyMs / Double(totalInferences)
    }

    public var totalCalls: Int { totalInferences }
}
