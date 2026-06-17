import SwiftUI

struct FeedbackView: View {
    @Environment(AppState.self) private var appState
    @State private var description = ""
    @State private var category = "misclassification"
    @State private var submitted = false
    @State private var traceID = ""

    var body: some View {
        Form {
            Section(String(localized: "feedback.details")) {
                Picker(String(localized: "feedback.category"), selection: $category) {
                    Text("misclassification").tag("misclassification")
                    Text("bug").tag("bug")
                    Text("feature").tag("feature")
                    Text("appeal").tag("appeal")
                }
                TextField(String(localized: "feedback.description"), text: $description, axis: .vertical)
                    .lineLimit(3 ... 8)
            }
            if !traceID.isEmpty {
                Section(String(localized: "feedback.traceID")) {
                    Text(traceID).font(.caption.monospaced())
                    Button(String(localized: "feedback.copyTraceID")) {
                        UIPasteboard.general.string = traceID
                    }
                }
            }
            Section {
                Button(String(localized: "feedback.submit")) {
                    Task { await submit() }
                }
                .disabled(description.isEmpty)
            }
            if submitted {
                Text(String(localized: "feedback.thanks")).foregroundStyle(.green)
            }
        }
        .navigationTitle(String(localized: "help.feedback.submit"))
    }

    private func submit() async {
        let payload = FeedbackPayload(description: description, category: category, traceID: UUID().uuidString)
        do {
            let response: FeedbackResponse = try await APIClient.shared.request(APIEndpoint(
                path: "/api/v1/feedback",
                method: "POST",
                body: try JSONEncoder().encode(payload)
            ))
            traceID = response.traceID
            appState.lastTraceID = response.traceID
            submitted = true
            AnalyticsManager.shared.track(.feedbackSubmitted(traceID: response.traceID))
        } catch {
            traceID = await APIClient.shared.lastTraceID
            appState.lastTraceID = traceID
            ErrorPresenter.shared.present(error)
        }
    }
}
