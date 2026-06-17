import CallKit
import Foundation
import SharedModels

final class CallDirectoryHandler: CXCallDirectoryProvider {
    private static let updateNotification = "com.msgguard.callblocklist.updated" as CFString

    override init() {
        super.init()
        registerForBlocklistUpdates()
    }

    override func beginRequest(with context: CXCallDirectoryExtensionContext) {
        let numbers = Self.loadBlockingNumbers()
        for number in numbers.sorted() {
            context.addBlockingEntry(withNextSequentialPhoneNumber: number)
        }
        context.completeRequest()
    }

    private func registerForBlocklistUpdates() {
        let center = CFNotificationCenterGetDarwinNotifyCenter()
        CFNotificationCenterAddObserver(
            center,
            Unmanaged.passUnretained(self).toOpaque(),
            { _, _, _, _, _ in
                CXCallDirectoryManager.sharedInstance.reloadExtension(
                    withIdentifier: "com.ethanshen.msgguard.CallDirectoryExtension",
                    completionHandler: nil
                )
            },
            Self.updateNotification,
            nil,
            .deliverImmediately
        )
    }

    private static func loadBlockingNumbers() -> [CXCallDirectoryPhoneNumber] {
        guard let url = FileManager.default.containerURL(forSecurityApplicationGroupIdentifier: AppConstants.appGroupID)?
            .appendingPathComponent(AppConstants.AppGroupFiles.callBlocklist),
              let data = try? Data(contentsOf: url),
              let decoded = try? JSONDecoder().decode([Int64].self, from: data) else {
            return []
        }
        return decoded
    }
}
