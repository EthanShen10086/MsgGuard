import CoreML
import Foundation
import SharedModels

public final class CoreMLClassifier: @unchecked Sendable {
    private let lock = NSLock()
    private var model: MLModel?
    private var featurizer: CoreMLFeaturizer?

    public init() {}

    public func load(from compiledURL: URL) throws {
        lock.lock()
        defer { lock.unlock() }
        model = try MLModel(contentsOf: compiledURL)
        let featurizerURL = compiledURL.deletingLastPathComponent().appendingPathComponent("coreml_featurizer.json")
        if FileManager.default.fileExists(atPath: featurizerURL.path) {
            featurizer = CoreMLFeaturizer(url: featurizerURL)
        } else {
            featurizer = nil
        }
    }

    public func loadFeaturizer(from url: URL) {
        lock.lock()
        defer { lock.unlock() }
        featurizer = CoreMLFeaturizer(url: url)
    }

    public func unload() {
        lock.lock()
        defer { lock.unlock() }
        model = nil
        featurizer = nil
    }

    public var isLoaded: Bool {
        lock.lock()
        defer { lock.unlock() }
        return model != nil
    }

    public func classify(text: String) -> FilterResult? {
        lock.lock()
        let active = model
        let localFeaturizer = featurizer
        lock.unlock()
        guard let active else { return nil }

        if let localFeaturizer, let features = localFeaturizer.vectorize(text) {
            guard let input = try? MLDictionaryFeatureProvider(dictionary: ["features": MLFeatureValue(multiArray: features)]),
                  let out = try? active.prediction(from: input) else {
                return nil
            }
            let score: Double
            if let value = out.featureValue(for: "spam_score")?.multiArrayValue, value.count > 0 {
                score = value[0].doubleValue
            } else if let value = out.featureValue(for: "spam_score")?.doubleValue {
                score = value
            } else {
                return nil
            }
            let isSpam = score >= localFeaturizer.threshold
            return FilterResult(
                category: isSpam ? .spam : .ham,
                confidence: isSpam ? score : 1.0 - score,
                layer: .coreML,
                shouldFilter: isSpam
            )
        }

        guard let input = try? MLDictionaryFeatureProvider(dictionary: ["text": MLFeatureValue(string: text)]),
              let out = try? active.prediction(from: input) else {
            return nil
        }

        let label: String
        if let s = out.featureValue(for: "label")?.stringValue {
            label = s
        } else if let n = out.featureValue(for: "label")?.int64Value {
            label = n == 1 ? "spam" : "ham"
        } else if let score = out.featureValue(for: "spam_score")?.doubleValue {
            let isSpam = score >= 0.5
            return FilterResult(
                category: isSpam ? .spam : .ham,
                confidence: isSpam ? score : 1.0 - score,
                layer: .coreML,
                shouldFilter: isSpam
            )
        } else {
            return nil
        }

        let isSpam = label == "spam" || label == "1" || label == "phishing" || label == "promotion"
        return FilterResult(
            category: isSpam ? .spam : .ham,
            confidence: isSpam ? 0.88 : 0.85,
            layer: .coreML,
            shouldFilter: isSpam
        )
    }
}

public enum CoreMLModelCompiler {
    public static func compileAndInstall(rawModelURL: URL, destinationURL: URL, featurizerURL: URL? = nil) throws {
        let compiled = try MLModel.compileModel(at: rawModelURL)
        let fm = FileManager.default
        if fm.fileExists(atPath: destinationURL.path) {
            try fm.removeItem(at: destinationURL)
        }
        try fm.createDirectory(at: destinationURL.deletingLastPathComponent(), withIntermediateDirectories: true)
        try fm.copyItem(at: compiled, to: destinationURL)
        if let featurizerURL, fm.fileExists(atPath: featurizerURL.path) {
            let destFeaturizer = destinationURL.deletingLastPathComponent().appendingPathComponent("coreml_featurizer.json")
            if fm.fileExists(atPath: destFeaturizer.path) {
                try fm.removeItem(at: destFeaturizer)
            }
            try fm.copyItem(at: featurizerURL, to: destFeaturizer)
        }
    }
}
