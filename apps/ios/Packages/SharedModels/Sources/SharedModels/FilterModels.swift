import Foundation

public enum MessageCategory: String, Codable, CaseIterable, Sendable {
    case ham
    case spam
    case promotion
    case transaction
    case finance
    case operatorMsg = "operator"
    case phishing

    public var displayKey: String {
        "category.\(rawValue)"
    }
}

public enum FilterTag: String, Codable, CaseIterable, Identifiable, Sendable {
    case advertising
    case charity
    case shortCode
    case finance
    case order
    case operatorMsg = "operator"

    public var id: String { rawValue }
}

public enum RuleType: String, Codable, Sendable {
    case numberAllow
    case numberBlock
    case keywordAllow
    case keywordBlock
}

public struct BlockRule: Codable, Identifiable, Sendable, Equatable {
    public let id: UUID
    public var type: RuleType
    public var pattern: String
    public var priority: Int
    public var enabled: Bool

    public init(
        id: UUID = UUID(),
        type: RuleType,
        pattern: String,
        priority: Int = 0,
        enabled: Bool = true
    ) {
        self.id = id
        self.type = type
        self.pattern = pattern
        self.priority = priority
        self.enabled = enabled
    }
}

public struct FilterConfig: Codable, Sendable, Equatable {
    public var enabledTags: Set<FilterTag>
    public var rules: [BlockRule]
    public var cloudLLMEnabled: Bool
    public var otpProtectionEnabled: Bool
    public var locale: String
    public var modelVersion: String
    public var rulesVersion: String
    public var rulesChecksum: String
    public var iCloudSyncEnabled: Bool

    public init(
        enabledTags: Set<FilterTag> = Set(FilterTag.allCases),
        rules: [BlockRule] = [],
        cloudLLMEnabled: Bool = false,
        otpProtectionEnabled: Bool = true,
        locale: String = "zh-Hans",
        modelVersion: String = "1.0.0",
        rulesVersion: String = "1.0.0",
        rulesChecksum: String = "",
        iCloudSyncEnabled: Bool = false
    ) {
        self.enabledTags = enabledTags
        self.rules = rules
        self.cloudLLMEnabled = cloudLLMEnabled
        self.otpProtectionEnabled = otpProtectionEnabled
        self.locale = locale
        self.modelVersion = modelVersion
        self.rulesVersion = rulesVersion
        self.rulesChecksum = rulesChecksum
        self.iCloudSyncEnabled = iCloudSyncEnabled
    }
}

public struct FilterResult: Sendable, Equatable {
    public let category: MessageCategory
    public let confidence: Double
    public let layer: FilterLayer
    public let shouldFilter: Bool

    public init(category: MessageCategory, confidence: Double, layer: FilterLayer, shouldFilter: Bool) {
        self.category = category
        self.confidence = confidence
        self.layer = layer
        self.shouldFilter = shouldFilter
    }
}

public enum FilterLayer: String, Codable, Sendable {
    case heuristic
    case naiveBayes
    case coreML
    case cloudLLM
    case rule
}

public struct FilterStats: Codable, Sendable, Equatable {
    public var blockedToday: Int
    public var blockedTotal: Int
    public var byCategory: [String: Int]
    public var lastUpdated: Date

    public init(
        blockedToday: Int = 0,
        blockedTotal: Int = 0,
        byCategory: [String: Int] = [:],
        lastUpdated: Date = Date()
    ) {
        self.blockedToday = blockedToday
        self.blockedTotal = blockedTotal
        self.byCategory = byCategory
        self.lastUpdated = lastUpdated
    }
}

public struct SampleEntry: Codable, Identifiable, Sendable {
    public let id: UUID
    public let text: String
    public let label: MessageCategory
    public let createdAt: Date

    public init(id: UUID = UUID(), text: String, label: MessageCategory, createdAt: Date = Date()) {
        self.id = id
        self.text = text
        self.label = label
        self.createdAt = createdAt
    }
}

public enum UserMode: String, Codable, Sendable, CaseIterable {
    case standard
    case elder
}
