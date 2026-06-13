import FilterEngine
import Foundation
import MailKit
import SharedModels

/// macOS Mail.app extension — classifies messages with HybridFilterEngine (subject + full body).
final class MailExtensionHandler: NSObject, MEExtension, MEMessageActionHandler {
    func handlerForMessageActions() -> MEMessageActionHandler { self }

    func decideAction(for message: MEMessage, completionHandler: @escaping (MEMessageActionDecision?) -> Void) {
        if message.rawData == nil {
            completionHandler(MEMessageActionDecision.invokeAgainWithBody)
            return
        }

        let config = SyncConfigLoader.loadConfig()
        var engine = HybridFilterEngine()
        if let bayes = SyncConfigLoader.loadBayesModel() {
            engine.loadBayesModel(from: bayes)
        }
        if let coreURL = SyncConfigLoader.coreMLModelURL(),
           FileManager.default.fileExists(atPath: coreURL.path) {
            try? engine.loadCoreML(from: coreURL)
        }

        let subject = message.subject
        let sender = message.fromAddress.addressString
        let bodyText = MailRFC822Parser.plainText(from: message.rawData) ?? ""
        let combined = [subject, bodyText].filter { !$0.isEmpty }.joined(separator: "\n")

        let result = engine.classify(sender: sender, body: combined, config: config)
        guard result.shouldFilter else {
            completionHandler(nil)
            return
        }
        completionHandler(MEMessageActionDecision.action(MEMessageAction.moveToTrash))
    }
}
