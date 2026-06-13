import Foundation
import SharedModels

public struct HybridFilterEngine: Sendable {
    private let heuristic = HeuristicFilter()
    private let bayes = NaiveBayesClassifier()
    private let coreML = CoreMLClassifier()
    public var healthMonitor = MLModelHealthMonitor()

    public init() {}

    public mutating func loadBayesModel(from data: Data) {
        try? bayes.importModel(from: data)
    }

    public mutating func classify(sender: String?, body: String, config: FilterConfig) -> FilterResult {
        let start = CFAbsoluteTimeGetCurrent()
        defer { recordLatency(start) }

        if let result = heuristic.classify(sender: sender, body: body, config: config),
           result.confidence >= 0.85 {
            return result
        }

        if let result = bayes.classify(text: body), result.confidence >= 0.7 {
            return result
        }

        if let result = coreML.classify(text: body) {
            return result
        }

        if let result = heuristic.classify(sender: sender, body: body, config: config) {
            return result
        }

        if let result = bayes.classify(text: body) {
            return result
        }

        return FilterResult(category: .ham, confidence: 0.5, layer: .heuristic, shouldFilter: false)
    }

    public func trainSample(text: String, category: MessageCategory) {
        bayes.train(text: text, category: category)
    }

    public func exportBayesModel() -> Data? {
        bayes.exportModel()
    }

    private mutating func recordLatency(_ start: CFAbsoluteTime) {
        let ms = (CFAbsoluteTimeGetCurrent() - start) * 1000
        healthMonitor.record(latencyMs: ms)
    }
}
