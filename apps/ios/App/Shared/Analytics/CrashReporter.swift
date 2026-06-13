import Foundation
import SharedModels

/// Opt-in crash report collector (JSONL to App Group).
final class CrashReporter: @unchecked Sendable {
    static let shared = CrashReporter()
    private let queue = DispatchQueue(label: "com.msgguard.crashreporter")
    private var enabled = false
    private var installed = false

    private init() {}

    func install() {
        guard !installed else { return }
        installed = true
        writeSentinel()
        NSSetUncaughtExceptionHandler { exception in
            CrashReporter.shared.recordException(exception)
        }
    }

    func setEnabled(_ value: Bool) {
        enabled = value
    }

    func record(error: Error, context: [String: String] = [:]) {
        guard enabled else { return }
        recordPayload(["error": String(describing: error), "context": context])
    }

    private func recordException(_ exception: NSException) {
        guard enabled else { return }
        recordPayload(["exception": exception.name.rawValue, "reason": exception.reason ?? ""])
    }

    private func recordPayload(_ fields: [String: Any]) {
        queue.async {
            var entry = fields
            entry["timestamp"] = ISO8601DateFormatter().string(from: Date())
            guard let data = try? JSONSerialization.data(withJSONObject: entry),
                  let line = String(data: data, encoding: .utf8) else { return }
            self.append(line + "\n")
        }
    }

    private func writeSentinel() {
        guard let url = FileManager.default.containerURL(
            forSecurityApplicationGroupIdentifier: AppConstants.appGroupID
        )?.appendingPathComponent("crash_reporter.installed") else { return }
        try? "ok".write(to: url, atomically: true, encoding: .utf8)
    }

    private func append(_ line: String) {
        guard let url = FileManager.default.containerURL(
            forSecurityApplicationGroupIdentifier: AppConstants.appGroupID
        )?.appendingPathComponent("crash_reports.jsonl") else { return }
        if FileManager.default.fileExists(atPath: url.path),
           let handle = try? FileHandle(forWritingTo: url) {
            handle.seekToEndOfFile()
            handle.write(line.data(using: .utf8)!)
            try? handle.close()
        } else {
            try? line.write(to: url, atomically: true, encoding: .utf8)
        }
        MGLogger.app.info("crash recorded")
    }
}
