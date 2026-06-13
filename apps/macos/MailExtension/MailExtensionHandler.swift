import FilterEngine
import Foundation
import MailKit
import SharedModels

/// macOS Mail.app extension — classifies message subject/body with the same HybridFilterEngine as SMS.
final class MailExtensionHandler: NSObject, MEExtension, MEMessageActionHandler {
    func handlerForMessageActions() -> MEMessageActionHandler { self }

    func decideAction(for message: MEMessage, completionHandler: @escaping (MEMessageActionDecision?) -> Void) {
        let config = SyncConfigLoader.loadConfig()
        var engine = HybridFilterEngine()
        if let bayes = SyncConfigLoader.loadBayesModel() {
            engine.loadBayesModel(from: bayes)
        }
        if let coreURL = SyncConfigLoader.coreMLModelURL(),
           FileManager.default.fileExists(atPath: coreURL.path) {
            try? engine.loadCoreML(from: coreURL)
        }

        let subject = message.subject ?? ""
        let body = subject
        let result = engine.classify(sender: nil, body: body, config: config)

        guard result.shouldFilter else {
            completionHandler(nil)
            return
        }
        completionHandler(MEMessageActionDecision.action(MEMessageAction.moveToTrash))
    }
}
