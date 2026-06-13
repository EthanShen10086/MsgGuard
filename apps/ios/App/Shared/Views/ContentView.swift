import DesignSystem
import SwiftUI

struct ContentView: View {
    @Environment(AppState.self) private var appState

    var body: some View {
        Group {
            if !appState.onboardingCompleted {
                OnboardingView()
            } else {
                MainTabView()
            }
        }
        .userMode(appState.userMode)
    }
}

struct MainTabView: View {
    var body: some View {
        TabView {
            DashboardView()
                .tabItem { Label(String(localized: "tab.home"), systemImage: "shield.checkered") }
            RulesView()
                .tabItem { Label(String(localized: "tab.rules"), systemImage: "list.bullet") }
            SamplesView()
                .tabItem { Label(String(localized: "tab.samples"), systemImage: "text.badge.plus") }
            StatsView()
                .tabItem { Label(String(localized: "tab.stats"), systemImage: "chart.bar") }
            HelpView()
                .tabItem { Label(String(localized: "tab.help"), systemImage: "questionmark.circle") }
            SettingsView()
                .tabItem { Label(String(localized: "tab.settings"), systemImage: "gearshape") }
        }
    }
}
