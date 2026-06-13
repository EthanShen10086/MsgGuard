import XCTest
@testable import FilterEngine
import SharedModels

final class FilterEngineTests: XCTestCase {
    func testHeuristicDetectsChineseSpam() {
        let filter = HeuristicFilter()
        let config = FilterConfig(locale: "zh-Hans")
        let result = filter.classify(sender: "+8613800138000", body: "恭喜中奖免费领取加微信", config: config)
        XCTAssertNotNil(result)
        XCTAssertTrue(result?.shouldFilter == true)
    }

    func testHeuristicAllowsVerificationCode() {
        let filter = HeuristicFilter()
        let config = FilterConfig(enabledTags: [], locale: "zh-Hans")
        let result = filter.classify(sender: "1008610086", body: "您的验证码是123456", config: config)
        XCTAssertEqual(result?.shouldFilter, false)
    }

    func testNaiveBayesClassifiesSpam() {
        let bayes = NaiveBayesClassifier()
        let result = bayes.classify(text: "免费贷款无抵押当天放款加微信")
        XCTAssertNotNil(result)
        XCTAssertTrue(result?.category == .spam || result?.category == .promotion)
    }

    func testHybridEngineEndToEnd() {
        var engine = HybridFilterEngine()
        let config = FilterConfig(locale: "zh-Hans")
        let result = engine.classify(sender: nil, body: "免费贷款无抵押当天放款", config: config)
        XCTAssertTrue(result.shouldFilter)
    }
}
