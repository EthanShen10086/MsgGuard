import Foundation

public enum AppConstants {
    public static let appGroupID = "group.com.ethanshen.msgguard"
    public static let bundleID = "com.ethanshen.msgguard"

    public enum UserDefaultsKeys {
        public static let extensionEnabled = "extensionEnabled"
        public static let onboardingCompleted = "onboardingCompleted"
        public static let blockedCountToday = "blockedCountToday"
        public static let blockedCountTotal = "blockedCountTotal"
        public static let lastBlockedDate = "lastBlockedDate"
        public static let cloudLLMEnabled = "cloudLLMEnabled"
        public static let preferredLocale = "preferredLocale"
        public static let userMode = "userMode"
        public static let lastTraceID = "lastTraceID"
    }

    public enum AppGroupFiles {
        public static let rulesManifest = "rules_manifest.json"
        public static let bayesModel = "bayes_model.json"
        public static let statsSnapshot = "stats_snapshot.json"
        public static let coreMLModel = "spam_classifier.mlmodelc"
    }

    public enum URLSchemes {
        public static let base = "msgguard"
        public static let dashboard = "msgguard://dashboard"
        public static let settings = "msgguard://settings"
    }
}
