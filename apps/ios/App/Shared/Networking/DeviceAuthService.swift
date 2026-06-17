import Foundation
import OSLog
import Security

/// Persists device ID and gateway device token in Keychain.
actor DeviceAuthService {
    static let shared = DeviceAuthService()

    private let logger = Logger(subsystem: "com.ethanshen.msgguard", category: "device-auth")
    private let tokenKey = "com.msgguard.device_token"
    private let deviceIDKey = "com.msgguard.device_id"
    private let baseURL: URL

    private init() {
        baseURL = URL(string: ProcessInfo.processInfo.environment["MSGGUARD_API_BASE"] ?? "http://localhost:8080")!
    }

    var deviceID: String {
        if let stored = readKeychain(account: deviceIDKey) { return stored }
        let id = UIDeviceIdentifier.current
        writeKeychain(account: deviceIDKey, value: id)
        return id
    }

    /// Ensures a valid device bearer token is available (POST /api/v1/auth/device).
    func ensureAuthenticated() async throws -> String {
        if let token = readKeychain(account: tokenKey), !token.isEmpty {
            return token
        }
        return try await fetchDeviceToken()
    }

    func fetchDeviceToken() async throws -> String {
        var request = URLRequest(url: baseURL.appendingPathComponent("/api/v1/auth/device"))
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        let body = ["device_id": deviceID]
        request.httpBody = try JSONSerialization.data(withJSONObject: body)

        let (data, response) = try await URLSession.shared.data(for: request)
        guard let http = response as? HTTPURLResponse, (200 ... 299).contains(http.statusCode) else {
            throw MGError.network(.serverError)
        }
        struct TokenResponse: Decodable {
            let access_token: String
        }
        let decoded = try JSONDecoder().decode(TokenResponse.self, from: data)
        writeKeychain(account: tokenKey, value: decoded.access_token)
        logger.info("Device token acquired")
        return decoded.access_token
    }

    func clearToken() {
        deleteKeychain(account: tokenKey)
    }

    private func readKeychain(account: String) -> String? {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrAccount as String: account,
            kSecReturnData as String: true,
            kSecMatchLimit as String: kSecMatchLimitOne,
        ]
        var result: AnyObject?
        guard SecItemCopyMatching(query as CFDictionary, &result) == errSecSuccess,
              let data = result as? Data,
              let value = String(data: data, encoding: .utf8) else {
            return nil
        }
        return value
    }

    private func writeKeychain(account: String, value: String) {
        deleteKeychain(account: account)
        guard let data = value.data(using: .utf8) else { return }
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrAccount as String: account,
            kSecValueData as String: data,
            kSecAttrAccessible as String: kSecAttrAccessibleAfterFirstUnlock,
        ]
        SecItemAdd(query as CFDictionary, nil)
    }

    private func deleteKeychain(account: String) {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrAccount as String: account,
        ]
        SecItemDelete(query as CFDictionary)
    }
}
