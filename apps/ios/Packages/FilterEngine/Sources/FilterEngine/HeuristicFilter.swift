import Foundation
import SharedModels

public struct HeuristicFilter: Sendable {
    private let spamKeywordsZH: [String] = [
        "免费领取", "中奖", "贷款", "套现", "刷单", "兼职", "推广", "优惠", "点击链接",
        "退订回T", "回复TD", "验证码请勿", "低息", "无抵押", "加微信", "加QQ",
    ]
    private let spamKeywordsEN: [String] = [
        "free gift", "winner", "click here", "unsubscribe", "loan", "casino",
        "verify account", "act now", "limited offer", "congratulations",
    ]
    private let hamKeywords: [String] = [
        "验证码", "verification code", "取件码", "快递", "delivery", "您的订单",
        "银行", "bank", "登录", "login",
    ]

    public init() {}

    public func classify(sender: String?, body: String, config: FilterConfig) -> FilterResult? {
        let text = body.lowercased()
        let senderDigits = sender?.filter(\.isNumber) ?? ""

        if matchesRules(text: text, sender: sender ?? "", config: config) {
            return FilterResult(category: .spam, confidence: 1.0, layer: .rule, shouldFilter: true)
        }

        if isShortCode(senderDigits) && config.enabledTags.contains(.shortCode) {
            return nil
        }

        for keyword in hamKeywords where text.contains(keyword.lowercased()) {
            return FilterResult(category: .ham, confidence: 0.95, layer: .heuristic, shouldFilter: false)
        }

        let keywords = config.locale.hasPrefix("zh") ? spamKeywordsZH : spamKeywordsEN
        var hits = 0
        for keyword in keywords where text.contains(keyword.lowercased()) {
            hits += 1
        }

        if hits >= 2 {
            return FilterResult(category: .spam, confidence: min(0.99, 0.7 + Double(hits) * 0.1), layer: .heuristic, shouldFilter: true)
        }
        if hits == 1 {
            return FilterResult(category: .promotion, confidence: 0.75, layer: .heuristic, shouldFilter: config.enabledTags.contains(.advertising))
        }

        if text.contains("http://") || text.contains("https://") || text.contains("bit.ly") {
            return FilterResult(category: .phishing, confidence: 0.8, layer: .heuristic, shouldFilter: true)
        }

        return nil
    }

    private func matchesRules(text: String, sender: String, config: FilterConfig) -> Bool {
        let sorted = config.rules.filter(\.enabled).sorted { $0.priority > $1.priority }
        for rule in sorted {
            switch rule.type {
            case .numberAllow:
                if sender.contains(rule.pattern) { return false }
            case .numberBlock:
                if sender.contains(rule.pattern) { return true }
            case .keywordAllow:
                if text.localizedCaseInsensitiveContains(rule.pattern) { return false }
            case .keywordBlock:
                if text.localizedCaseInsensitiveContains(rule.pattern) { return true }
            }
        }
        return false
    }

    private func isShortCode(_ digits: String) -> Bool {
        digits.count == 5 || digits.hasPrefix("106")
    }
}
