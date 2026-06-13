import Foundation
import SharedModels

public struct CoreMLClassifier: Sendable {
    public init() {}

    public func classify(text: String) -> FilterResult? {
        // L2 placeholder — model loaded from App Group after Phase 2 training pipeline.
        // Returns nil to fall through to next layer when no bundled model is present.
        nil
    }
}

public struct MLModelHealthMonitor: Sendable {
    public private(set) var totalInferences = 0
    public private(set) var totalLatencyMs: Double = 0
    public private(set) var lastLatencyMs: Double = 0

    public init() {}

    public mutating func record(latencyMs: Double) {
        totalInferences += 1
        totalLatencyMs += latencyMs
        lastLatencyMs = latencyMs
    }

    public var meanLatencyMs: Double {
        guard totalInferences > 0 else { return 0 }
        return totalLatencyMs / Double(totalInferences)
    }
}
