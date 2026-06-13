import Foundation
import SharedModels

public final class NaiveBayesClassifier: @unchecked Sendable {
    private var wordCounts: [MessageCategory: [String: Int]] = [:]
    private var categoryCounts: [MessageCategory: Int] = [:]
    private var totalDocuments = 0
    private let smoothing = 1.0

    public init() {
        seedDefaults()
    }

    public func train(text: String, category: MessageCategory) {
        let tokens = tokenize(text)
        totalDocuments += 1
        categoryCounts[category, default: 0] += 1
        for token in tokens {
            var counts = wordCounts[category, default: [:]]
            counts[token, default: 0] += 1
            wordCounts[category] = counts
        }
    }

    public func classify(text: String) -> FilterResult? {
        guard totalDocuments > 0 else { return nil }
        let tokens = tokenize(text)
        guard !tokens.isEmpty else { return nil }

        var bestCategory = MessageCategory.ham
        var bestScore = -Double.infinity

        for category in MessageCategory.allCases {
            let count = categoryCounts[category, default: 0]
            guard count > 0 else { continue }
            var score = log(Double(count) / Double(totalDocuments))
            let vocabSize = wordCounts.values.reduce(0) { $0 + $1.count }
            for token in tokens {
                let wordCount = wordCounts[category]?[token] ?? 0
                let categoryTotal = wordCounts[category]?.values.reduce(0, +) ?? 0
                let prob = (Double(wordCount) + smoothing) / (Double(categoryTotal) + smoothing * Double(vocabSize + 1))
                score += log(prob)
            }
            if score > bestScore {
                bestScore = score
                bestCategory = category
            }
        }

        let confidence = min(0.99, max(0.5, 1.0 / (1.0 + exp(-bestScore / Double(tokens.count)))))
        let shouldFilter = bestCategory == .spam || bestCategory == .promotion || bestCategory == .phishing
        guard confidence >= 0.5 else { return nil }

        return FilterResult(
            category: bestCategory,
            confidence: confidence,
            layer: .naiveBayes,
            shouldFilter: shouldFilter
        )
    }

    public func exportModel() -> Data? {
        let payload = BayesModelPayload(
            wordCounts: wordCounts.mapKeys { $0.rawValue },
            categoryCounts: categoryCounts.mapKeys { $0.rawValue },
            totalDocuments: totalDocuments
        )
        return try? JSONEncoder().encode(payload)
    }

    public func importModel(from data: Data) throws {
        let payload = try JSONDecoder().decode(BayesModelPayload.self, from: data)
        wordCounts = payload.wordCounts.compactMapKeys { MessageCategory(rawValue: $0) }
        categoryCounts = payload.categoryCounts.compactMapKeys { MessageCategory(rawValue: $0) }
        totalDocuments = payload.totalDocuments
    }

    private func tokenize(_ text: String) -> [String] {
        text.lowercased()
            .components(separatedBy: CharacterSet.alphanumerics.inverted)
            .filter { $0.count > 1 }
    }

    private func seedDefaults() {
        let spamSamples = [
            "恭喜您中奖了请点击链接领取奖品",
            "免费贷款无抵押当天放款加微信",
            "推广优惠活动限时抢购回复TD退订",
            "Congratulations you won click here now",
            "Free gift claim your prize unsubscribe",
        ]
        let hamSamples = [
            "您的验证码是123456请勿泄露",
            "快递已到达请凭取件码领取",
            "Your verification code is 654321",
            "Your order has been shipped",
        ]
        for sample in spamSamples { train(text: sample, category: .spam) }
        for sample in hamSamples { train(text: sample, category: .ham) }
    }
}

private struct BayesModelPayload: Codable {
    let wordCounts: [String: [String: Int]]
    let categoryCounts: [String: Int]
    let totalDocuments: Int
}

private extension Dictionary {
    func mapKeys<T: Hashable>(_ transform: (Key) -> T) -> [T: Value] {
        var result: [T: Value] = [:]
        for (key, value) in self {
            result[transform(key)] = value
        }
        return result
    }

    func compactMapKeys<T: Hashable>(_ transform: (Key) -> T?) -> [T: Value] {
        var result: [T: Value] = [:]
        for (key, value) in self {
            if let newKey = transform(key) {
                result[newKey] = value
            }
        }
        return result
    }
}
