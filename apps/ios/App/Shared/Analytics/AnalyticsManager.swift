import Foundation
import SharedModels

public enum AnalyticsEvent: Sendable {
    case appLaunched
    case onboardingCompleted
    case extensionEnabledChanged(enabled: Bool)
    case filterStarted
    case filterCompleted(category: String, layer: String)
    case sampleSubmitted(label: String)
    case feedbackSubmitted(traceID: String)
    case settingsChanged(key: String, value: String)
    case error(domain: String, code: String)
    case purchaseCompleted(productId: String)
}

public protocol AnalyticsTracking: Sendable {
    func track(_ event: AnalyticsEvent)
}

public final class AnalyticsManager: AnalyticsTracking, @unchecked Sendable {
    public static let shared = AnalyticsManager()
    private init() {}

    public func track(_ event: AnalyticsEvent) {
        MGLogger.app.info("analytics: \(String(describing: event))")
    }
}
