import Foundation

enum MGError: Error, LocalizedError, Identifiable {
    case filter(FilterFailure)
    case store(StoreFailure)
    case network(NetworkFailure)
    case generic(String)

    var id: String { "\(analyticsDomain).\(analyticsCode).\(UUID().uuidString)" }

    enum FilterFailure: String {
        case extensionNotEnabled
        case configLoadFailed
        case ruleSaveFailed
    }

    enum StoreFailure: String {
        case containerUnavailable
        case encodingFailed
    }

    enum NetworkFailure: String {
        case offline
        case timeout
        case serverError
    }

    var errorDescription: String? {
        switch self {
        case let .filter(f): "Filter error: \(f.rawValue)"
        case let .store(f): "Store error: \(f.rawValue)"
        case let .network(f): "Network error: \(f.rawValue)"
        case let .generic(msg): msg
        }
    }

    var userMessage: String {
        switch self {
        case let .filter(f):
            switch f {
            case .extensionNotEnabled: String(localized: "error.filter.extension")
            case .configLoadFailed: String(localized: "error.filter.config")
            case .ruleSaveFailed: String(localized: "error.filter.save")
            }
        case let .store(f):
            switch f {
            case .containerUnavailable: String(localized: "error.store.container")
            case .encodingFailed: String(localized: "error.store.encoding")
            }
        case let .network(f):
            switch f {
            case .offline: String(localized: "error.network.offline")
            case .timeout: String(localized: "error.network.timeout")
            case .serverError: String(localized: "error.network.server")
            }
        case let .generic(msg): msg
        }
    }

    var analyticsDomain: String {
        switch self {
        case .filter: "filter"
        case .store: "store"
        case .network: "network"
        case .generic: "generic"
        }
    }

    var analyticsCode: String {
        switch self {
        case let .filter(f): f.rawValue
        case let .store(f): f.rawValue
        case let .network(f): f.rawValue
        case .generic: "generic"
        }
    }
}

@MainActor
final class ErrorPresenter {
    static let shared = ErrorPresenter()
    var currentError: MGError?

    private init() {}

    func present(_ error: Error) {
        let mgError = error as? MGError ?? .generic(error.localizedDescription)
        currentError = mgError
        AnalyticsManager.shared.track(.error(domain: mgError.analyticsDomain, code: mgError.analyticsCode))
    }
}
