import Foundation
import SharedModels

/// Protects OTP / verification / delivery codes from false-positive blocks.
public enum OTPGuard {
    private static let otpPatterns: [String] = [
        #"验证码[是为：:\s]*\d{4,8}"#,
        #"verification code[:\s]*\d{4,8}"#,
        #"取件码[是为：:\s]*\d{4,8}"#,
        #"pickup code[:\s]*\d{4,8}"#,
        #"动态码[是为：:\s]*\d{4,8}"#,
        #"登录码[是为：:\s]*\d{4,8}"#,
        #"【[\w\s]+】验证码"#,
    ]

    public static func isProtectedMessage(body: String, sender: String?) -> Bool {
        let lower = body.lowercased()
        if lower.contains("验证码") || lower.contains("verification code") {
            return true
        }
        if lower.contains("取件码") || lower.contains("pickup code") || lower.contains("快递") {
            return true
        }
        for pattern in otpPatterns {
            if body.range(of: pattern, options: .regularExpression) != nil {
                return true
            }
        }
        if isTrustedSender(sender) && lower.range(of: #"\d{4,8}"#, options: .regularExpression) != nil {
            return true
        }
        return false
    }

    private static func isTrustedSender(_ sender: String?) -> Bool {
        guard let sender else { return false }
        let digits = sender.filter(\.isNumber)
        return digits.hasPrefix("106") || digits.count <= 6
    }
}
