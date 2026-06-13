import FilterEngine
import IdentityLookup
import SharedModels

final class MessageFilterExtension: ILMessageFilterExtension {}

extension MessageFilterExtension: ILMessageFilterQueryHandling {
    func handle(
        _ queryRequest: ILMessageFilterQueryRequest,
        context: ILMessageFilterExtensionContext,
        completion: @escaping (ILMessageFilterQueryResponse) -> Void
    ) {
        let signpost = PerformanceMonitor.beginFilter()
        let response = ILMessageFilterQueryResponse()

        if queryRequest.messageBody == nil && queryRequest.sender == nil {
            response.action = .none
            completion(response)
            PerformanceMonitor.endFilter(signpost, latencyMs: 0)
            return
        }

        let config = SyncConfigLoader.loadConfig()
        if config.cloudLLMEnabled {
            context.deferQueryRequestToNetwork { networkResponse, error in
                if let networkResponse {
                    response.action = self.action(for: networkResponse)
                } else {
                    response.action = .none
                    MGLogger.network.error("defer failed: \(error?.localizedDescription ?? "unknown")")
                }
                completion(response)
            }
            return
        }

        var engine = HybridFilterEngine()
        if let modelData = SyncConfigLoader.loadBayesModel() {
            engine.loadBayesModel(from: modelData)
        }
        if let coreURL = SyncConfigLoader.coreMLModelURL(),
           FileManager.default.fileExists(atPath: coreURL.path) {
            engine.loadCoreML(from: coreURL)
        }

        let body = queryRequest.messageBody ?? ""
        let result = engine.classify(sender: queryRequest.sender, body: body, config: config)

        if result.shouldFilter {
            response.action = .junk
            SyncConfigLoader.saveStatsIncrement(category: result.category)
        } else {
            response.action = .allow
        }

        PerformanceMonitor.endFilter(signpost, latencyMs: engine.healthMonitor.meanLatencyMs)
        completion(response)
    }

    private func action(for networkResponse: ILNetworkResponse) -> ILMessageFilterAction {
        let data = networkResponse.data
        guard let json = try? JSONSerialization.jsonObject(with: data) as? [String: Any],
              let action = json["action"] as? String else {
            return .none
        }
        switch action {
        case "junk", "filter": return .junk
        case "promotion": return .promotion
        default: return .none
        }
    }
}
