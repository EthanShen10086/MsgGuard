import CallKit
import Foundation
import SharedModels

final class CallDirectoryHandler: CXCallDirectoryProvider {
    override func beginRequest(with context: CXCallDirectoryExtensionContext) {
        for number in Self.loadBlockingNumbers().sorted() {
            context.addBlockingEntry(withNextSequentialPhoneNumber: number)
        }
        context.completeRequest()
    }

    private static func loadBlockingNumbers() -> [CXCallDirectoryPhoneNumber] {
        guard let url = FileManager.default.containerURL(forSecurityApplicationGroupIdentifier: AppConstants.appGroupID)?
            .appendingPathComponent("call_blocklist.json"),
              let data = try? Data(contentsOf: url),
              let decoded = try? JSONDecoder().decode([Int64].self, from: data) else {
            return []
        }
        return decoded
    }
}
