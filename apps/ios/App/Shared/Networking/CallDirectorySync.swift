import CallKit
import Foundation
import SharedModels

/// Syncs call blocklist to App Group and reloads the Call Directory extension.
enum CallDirectorySync {
    private static let notificationName = "com.msgguard.callblocklist.updated" as CFString
    static let extensionBundleID = "com.ethanshen.msgguard.CallDirectoryExtension"

    static var blocklistURL: URL? {
        FileManager.default.containerURL(forSecurityApplicationGroupIdentifier: AppConstants.appGroupID)?
            .appendingPathComponent(AppConstants.AppGroupFiles.callBlocklist)
    }

    /// Writes numbers to App Group and triggers extension reload (OTA hook).
    @discardableResult
    static func applyOTAUpdate(numbers: [Int64]) throws -> URL {
        guard let url = blocklistURL else {
            throw MGError.store(.containerUnavailable)
        }
        let data = try JSONEncoder().encode(numbers)
        try data.write(to: url, options: .atomic)
        postUpdateNotification()
        reloadExtension()
        return url
    }

    static func reloadExtension() {
        CXCallDirectoryManager.sharedInstance.reloadExtension(withIdentifier: Self.extensionBundleID) { error in
            if let error {
                MGLogger.sync.error("CallDirectory reload failed: \(error.localizedDescription)")
            } else {
                MGLogger.sync.info("CallDirectory extension reloaded")
            }
        }
    }

    static func postUpdateNotification() {
        CFNotificationCenterPostNotification(
            CFNotificationCenterGetDarwinNotifyCenter(),
            CFNotificationName(notificationName),
            nil,
            nil,
            true
        )
    }
}
