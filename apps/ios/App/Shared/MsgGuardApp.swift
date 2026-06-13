import SwiftUI

@main
struct MsgGuardApp: App {
    @State private var appState = AppState()

    var body: some Scene {
        WindowGroup {
            ContentView()
                .environment(appState)
                .task {
                    await appState.loadAll()
                    AnalyticsManager.shared.track(.appLaunched)
                    let sync = SyncService()
                    try? await sync.syncRules()
                    try? await ModelUpdateService().checkAndUpdate()
                }
                .alert(
                    String(localized: "error.title"),
                    isPresented: Binding(
                        get: { ErrorPresenter.shared.currentError != nil },
                        set: { if !$0 { ErrorPresenter.shared.currentError = nil } }
                    )
                ) {
                    Button(String(localized: "error.ok")) {
                        ErrorPresenter.shared.currentError = nil
                    }
                } message: {
                    Text(ErrorPresenter.shared.currentError?.userMessage ?? "")
                }
        }
    }
}
