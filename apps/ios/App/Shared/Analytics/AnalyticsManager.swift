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
    case entitlementGranted(source: String)
    case entitlementRevoked
}

public protocol AnalyticsTracking: Sendable {
    func track(_ event: AnalyticsEvent)
}

public final class AnalyticsManager: AnalyticsTracking, @unchecked Sendable {
    public static let shared = AnalyticsManager()
    private let queue = DispatchQueue(label: "com.msgguard.analytics")
    private var uploadEnabled = true
    private var pending: [[String: Any]] = []
    private let client: NetworkClient = APIClient.shared

    private init() {}

    public func setUploadEnabled(_ enabled: Bool) {
        uploadEnabled = enabled
    }

    public func track(_ event: AnalyticsEvent) {
        let payload = eventPayload(event)
        MGLogger.app.info("analytics: \(payload["name"] as? String ?? "")")
        var shouldFlush = false
        queue.sync {
            appendJSONL(payload)
            pending.append(payload)
            shouldFlush = uploadEnabled && pending.count >= 5
        }
        if shouldFlush {
            Task(priority: .utility) { await self.flush() }
        }
    }

    public func flush() async {
        var batch: [[String: Any]] = []
        queue.sync {
            batch = pending
            pending = []
        }
        guard !batch.isEmpty else { return }
        for item in batch {
            guard let data = try? JSONSerialization.data(withJSONObject: item) else { continue }
            do {
                struct OK: Decodable { let status: String }
                let _: OK = try await client.request(APIEndpoint(path: "/api/v1/analytics", method: "POST", body: data))
            } catch {
                MGLogger.network.error("analytics upload failed")
            }
        }
    }

    private func eventPayload(_ event: AnalyticsEvent) -> [String: Any] {
        var props: [String: Any] = ["trace_id": TraceContext.lastTraceID]
        let name: String
        switch event {
        case .appLaunched: name = "app_launched"
        case .onboardingCompleted: name = "onboarding_completed"
        case .extensionEnabledChanged(let enabled):
            name = "extension_enabled_changed"; props["enabled"] = enabled
        case .filterStarted: name = "filter_started"
        case .filterCompleted(let category, let layer):
            name = "filter_completed"; props["category"] = category; props["layer"] = layer
        case .sampleSubmitted(let label): name = "sample_submitted"; props["label"] = label
        case .feedbackSubmitted(let traceID): name = "feedback_submitted"; props["trace_id"] = traceID
        case .settingsChanged(let key, let value): name = "settings_changed"; props["key"] = key; props["value"] = value
        case .error(let domain, let code): name = "error"; props["domain"] = domain; props["code"] = code
        case .purchaseCompleted(let productId): name = "purchase_completed"; props["product_id"] = productId
        case .entitlementGranted(let source): name = "entitlement_granted"; props["source"] = source
        case .entitlementRevoked: name = "entitlement_revoked"
        }
        return [
            "name": name,
            "props": props,
            "device_id": UIDeviceIdentifier.current,
            "timestamp": ISO8601DateFormatter().string(from: Date()),
        ]
    }

    private func appendJSONL(_ payload: [String: Any]) {
        guard let data = try? JSONSerialization.data(withJSONObject: payload),
              let line = String(data: data, encoding: .utf8) else { return }
        guard let url = FileManager.default.containerURL(
            forSecurityApplicationGroupIdentifier: AppConstants.appGroupID
        )?.appendingPathComponent("analytics.jsonl") else { return }
        if FileManager.default.fileExists(atPath: url.path),
           let handle = try? FileHandle(forWritingTo: url) {
            handle.seekToEndOfFile()
            handle.write((line + "\n").data(using: .utf8)!)
            try? handle.close()
        } else {
            try? (line + "\n").write(to: url, atomically: true, encoding: .utf8)
        }
    }
}

enum UIDeviceIdentifier {
    static var current: String {
        if let id = UserDefaults.standard.string(forKey: "mg_device_id") { return id }
        let id = UUID().uuidString
        UserDefaults.standard.set(id, forKey: "mg_device_id")
        return id
    }
}
