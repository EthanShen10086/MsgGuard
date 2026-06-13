import CoreML
import Foundation

struct CoreMLFeaturizer: Sendable {
    let featureCount: Int
    let vocabulary: [String: Int]
    let idf: [Double]
    let threshold: Double

    init?(url: URL) {
        guard let data = try? Data(contentsOf: url),
              let json = try? JSONSerialization.jsonObject(with: data) as? [String: Any],
              let featureCount = json["feature_count"] as? Int,
              let vocabularyRaw = json["vocabulary"] as? [String: Any],
              let idf = json["idf"] as? [Double] else {
            return nil
        }
        var vocabulary: [String: Int] = [:]
        for (term, value) in vocabularyRaw {
            if let idx = value as? Int {
                vocabulary[term] = idx
            } else if let idx = (value as? NSNumber)?.intValue {
                vocabulary[term] = idx
            }
        }
        self.featureCount = featureCount
        self.vocabulary = vocabulary
        self.idf = idf
        self.threshold = json["threshold"] as? Double ?? 0.5
    }

    func vectorize(_ text: String) -> MLMultiArray? {
        guard let array = try? MLMultiArray(shape: [NSNumber(value: featureCount)], dataType: .double) else {
            return nil
        }
        for i in 0 ..< featureCount {
            array[i] = 0
        }
        let tokens = text.lowercased()
            .components(separatedBy: CharacterSet.alphanumerics.inverted)
            .filter { $0.count > 1 }
        var counts: [String: Int] = [:]
        for token in tokens {
            counts[token, default: 0] += 1
        }
        for (token, count) in counts {
            guard let idx = vocabulary[token], idx < idf.count else { continue }
            array[idx] = NSNumber(value: Double(count) * idf[idx])
        }
        return array
    }
}
