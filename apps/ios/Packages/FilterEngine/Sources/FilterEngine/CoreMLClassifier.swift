import CoreML
import Foundation
import SharedModels

public final class CoreMLClassifier: @unchecked Sendable {
    private let lock = NSLock()
    private var model: MLModel?

    public init() {}

    public func load(from compiledURL: URL) throws {
        lock.lock()
        defer { lock.unlock() }
        model = try MLModel(contentsOf: compiledURL)
    }

    public func unload() {
        lock.lock()
        defer { lock.unlock() }
        model = nil
    }

    public var isLoaded: Bool {
        lock.lock()
        defer { lock.unlock() }
        return model != nil
    }

    public func classify(text: String) -> FilterResult? {
        lock.lock()
        let active = model
        lock.unlock()
        guard let active else { return nil }

        guard let input = try? MLDictionaryFeatureProvider(dictionary: ["text": MLFeatureValue(string: text)]),
              let out = try? active.prediction(from: input) else {
            return nil
        }

        let label: String
        if let s = out.featureValue(for: "label")?.stringValue {
            label = s
        } else if let n = out.featureValue(for: "label")?.int64Value {
            label = n == 1 ? "spam" : "ham"
        } else {
            return nil
        }

        let isSpam = label == "spam" || label == "1" || label == "phishing" || label == "promotion"
        let category: MessageCategory = isSpam ? .spam : .ham
        let confidence: Double = isSpam ? 0.88 : 0.85
        return FilterResult(
            category: category,
            confidence: confidence,
            layer: .coreML,
            shouldFilter: isSpam
        )
    }
}

public enum CoreMLModelCompiler {
    public static func compileAndInstall(rawModelURL: URL, destinationURL: URL) throws {
        let compiled = try MLModel.compileModel(at: rawModelURL)
        let fm = FileManager.default
        if fm.fileExists(atPath: destinationURL.path) {
            try fm.removeItem(at: destinationURL)
        }
        try fm.createDirectory(at: destinationURL.deletingLastPathComponent(), withIntermediateDirectories: true)
        try fm.copyItem(at: compiled, to: destinationURL)
    }
}
