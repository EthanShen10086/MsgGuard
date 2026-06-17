import Foundation

public struct APIEndpoint: Sendable {
    public let path: String
    public let method: String
    public let body: Data?

    public init(path: String, method: String = "GET", body: Data? = nil) {
        self.path = path
        self.method = method
        self.body = body
    }
}

public struct FeedbackPayload: Codable, Sendable {
    public let description: String
    public let category: String
    public let traceID: String
}

public struct FeedbackResponse: Codable, Sendable {
    public let id: String
    public let traceID: String
}

public struct ClassifyResponse: Codable, Sendable {
    public let action: String
    public let category: String
    public let confidence: Double
}

public enum APIFetchResult: Sendable {
    case data(Data)
    case notModified
}

public actor APIClient {
    public static let shared = APIClient()

    private let baseURL: URL
    private let session: URLSession
    public private(set) var lastTraceID = UUID().uuidString

    public init(baseURL: URL = URL(string: ProcessInfo.processInfo.environment["MSGGUARD_API_BASE"] ?? "http://localhost:8080")!) {
        self.baseURL = baseURL
        let config = URLSessionConfiguration.default
        config.timeoutIntervalForRequest = 15
        self.session = URLSession(configuration: config)
    }

    public func request<T: Decodable & Sendable>(_ endpoint: APIEndpoint) async throws -> T {
        let traceID = UUID().uuidString
        lastTraceID = traceID

        var request = URLRequest(url: baseURL.appendingPathComponent(endpoint.path))
        request.httpMethod = endpoint.method
        request.httpBody = endpoint.body
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue(traceID, forHTTPHeaderField: "X-Request-ID")
        if let token = try? await DeviceAuthService.shared.ensureAuthenticated() {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        let (data, response) = try await session.data(for: request)
        guard let http = response as? HTTPURLResponse else {
            throw MGError.network(.serverError)
        }
        guard (200 ... 299).contains(http.statusCode) else {
            throw MGError.network(.serverError)
        }
        if let returnedTrace = http.value(forHTTPHeaderField: "X-Request-ID") {
            lastTraceID = returnedTrace
            TraceContext.update(returnedTrace)
        }
        return try JSONDecoder().decode(T.self, from: data)
    }

    public func fetch(_ endpoint: APIEndpoint, ifNoneMatch: String? = nil) async throws -> APIFetchResult {
        let traceID = UUID().uuidString
        lastTraceID = traceID

        var request = URLRequest(url: baseURL.appendingPathComponent(endpoint.path))
        request.httpMethod = endpoint.method
        request.httpBody = endpoint.body
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue(traceID, forHTTPHeaderField: "X-Request-ID")
        if let token = try? await DeviceAuthService.shared.ensureAuthenticated() {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }
        if let etag = ifNoneMatch, !etag.isEmpty {
            request.setValue("\"\(etag)\"", forHTTPHeaderField: "If-None-Match")
        }

        let (data, response) = try await session.data(for: request)
        guard let http = response as? HTTPURLResponse else {
            throw MGError.network(.serverError)
        }
        if http.statusCode == 304 {
            return .notModified
        }
        guard (200 ... 299).contains(http.statusCode) else {
            throw MGError.network(.serverError)
        }
        if let returnedTrace = http.value(forHTTPHeaderField: "X-Request-ID") {
            lastTraceID = returnedTrace
            TraceContext.update(returnedTrace)
        }
        return .data(data)
    }
}
