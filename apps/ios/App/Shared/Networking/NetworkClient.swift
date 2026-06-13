import Foundation
import SharedModels

public protocol NetworkClient: Sendable {
    func request<T: Decodable & Sendable>(_ endpoint: APIEndpoint) async throws -> T
}

extension APIClient: NetworkClient {}

enum TraceContext: Sendable {
    private static let lock = NSLock()
    nonisolated(unsafe) private static var _lastTraceID = UUID().uuidString

    static var lastTraceID: String {
        lock.lock()
        defer { lock.unlock() }
        return _lastTraceID
    }

    static func update(_ id: String) {
        lock.lock()
        _lastTraceID = id
        lock.unlock()
    }
}
