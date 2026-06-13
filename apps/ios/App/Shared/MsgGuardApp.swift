import SharedModels
import SwiftUI

@main
struct MsgGuardApp: App {
    @State private var appState = AppState()

    init() {
        CrashReporter.shared.install()
    }

    var body: some Scene {
        WindowGroup {
            ContentView()
                .environment(appState)
                .task {
                    await appState.loadAll()
                    Task(priority: .utility) {
                        AnalyticsManager.shared.track(.appLaunched)
                        let sync = SyncService()
                        try? await sync.syncRules()
                        try? await ModelUpdateService().checkAndUpdate()
                        let perf = PerformanceMonitor.loadAggregateStats()
                        if !perf.isEmpty {
                            AnalyticsManager.shared.track(.filterCompleted(
                                category: "aggregate",
                                layer: "p99=\(Int(perf["max_ms"] ?? 0))ms"
                            ))
                        }
                        await AnalyticsManager.shared.flush()
                    }
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
