import BlocklistStore
import CryptoKit
import FilterEngine
import SharedModels
import WidgetKit

/// Downloads model artifacts from backend and writes to App Group.
actor ModelUpdateService {
    private let client = APIClient.shared
    private let store = BlocklistStore()

    struct ModelMeta: Decodable {
        let version: String
        let checksum: String
        let url: String
        let bayes_url: String?
    }

    func checkAndUpdate(locale: String? = nil) async throws {
        let config = try await store.loadConfig()
        let activeLocale = locale ?? config.locale
        let meta: ModelMeta = try await client.request(
            APIEndpoint(path: "/api/v1/models/latest?locale=\(activeLocale.addingPercentEncoding(withAllowedCharacters: .urlQueryAllowed) ?? activeLocale)")
        )
        if config.modelVersion == meta.version, await store.hasCoreMLModel(locale: activeLocale) {
            MGLogger.sync.info("Model up to date \(meta.version) locale=\(activeLocale)")
            return
        }
        MGLogger.sync.info("Updating model to \(meta.version) locale=\(activeLocale)")
        let artifactURL = resolveURL(meta.url)
        var request = URLRequest(url: artifactURL)
        request.httpMethod = "GET"
        let (data, response) = try await URLSession.shared.data(for: request)
        guard let http = response as? HTTPURLResponse, (200 ... 299).contains(http.statusCode) else {
            throw MGError.network(.serverError)
        }
        let digest = SHA256.hash(data: data).map { String(format: "%02x", $0) }.joined()
        let expected = meta.checksum.replacingOccurrences(of: "sha256:", with: "")
        if !expected.isEmpty, expected != "seed", digest != expected {
            MGLogger.sync.error("Model checksum mismatch")
            throw MGError.network(.serverError)
        }
        let rawURL = try await store.saveRawCoreMLModel(data, locale: activeLocale)
        let compiled = try await store.coreMLCompiledURL(locale: activeLocale)
        try CoreMLModelCompiler.compileAndInstall(rawModelURL: rawURL, destinationURL: compiled)
        try? FileManager.default.removeItem(at: rawURL)
        var updated = config
        updated.modelVersion = meta.version
        updated.locale = activeLocale
        try await store.saveConfig(updated)
        if let bayesPath = meta.bayes_url {
            try await downloadBayes(relativePath: bayesPath, locale: activeLocale)
        }
        WidgetCenter.shared.reloadAllTimelines()
    }

    private func downloadBayes(relativePath: String, locale: String) async throws {
        let url = resolveURL(relativePath)
        let (data, response) = try await URLSession.shared.data(for: URLRequest(url: url))
        guard let http = response as? HTTPURLResponse, (200 ... 299).contains(http.statusCode) else {
            return
        }
        try await store.saveBayesModel(data, locale: locale)
        MGLogger.sync.info("Bayes model updated locale=\(locale)")
    }

    private func resolveURL(_ path: String) -> URL {
        if path.hasPrefix("http") {
            return URL(string: path)!
        }
        let base = APIClient.sharedBaseURL
        return base.appendingPathComponent(path.trimmingCharacters(in: CharacterSet(charactersIn: "/")))
    }
}

extension APIClient {
    static var sharedBaseURL: URL {
        URL(string: ProcessInfo.processInfo.environment["MSGGUARD_API_BASE"] ?? "http://localhost:8080")!
    }
}
