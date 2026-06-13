import Foundation

enum MailRFC822Parser {
    /// Extracts plain text body from RFC822 raw message data when available.
    static func plainText(from rawData: Data?) -> String? {
        guard let rawData, !rawData.isEmpty,
              let raw = String(data: rawData, encoding: .utf8) ?? String(data: rawData, encoding: .isoLatin1) else {
            return nil
        }
        let parts = raw.components(separatedBy: "\r\n\r\n")
        guard parts.count >= 2 else { return nil }
        let headerBlock = parts[0]
        let body = parts.dropFirst().joined(separator: "\r\n\r\n")
        if headerBlock.lowercased().contains("content-type: text/html") {
            return stripHTML(body)
        }
        return body.trimmingCharacters(in: .whitespacesAndNewlines)
    }

    static func headerValue(_ name: String, in headers: [String: [String]]?) -> String? {
        guard let headers else { return nil }
        for (key, values) in headers where key.caseInsensitiveCompare(name) == .orderedSame {
            return values.first
        }
        return nil
    }

    private static func stripHTML(_ html: String) -> String {
        html.replacingOccurrences(of: "<[^>]+>", with: " ", options: .regularExpression)
            .replacingOccurrences(of: "\\s+", with: " ", options: .regularExpression)
            .trimmingCharacters(in: .whitespacesAndNewlines)
    }
}
