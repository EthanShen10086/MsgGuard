import Charts
import DesignSystem
import SwiftUI

struct StatsView: View {
    @Environment(AppState.self) private var appState

    var body: some View {
        NavigationStack {
            ScrollView {
                VStack(spacing: 16) {
                    MGCard {
                        VStack(alignment: .leading) {
                            Text(String(localized: "stats.overview"))
                                .font(.headline)
                            Text(String(localized: "dashboard.blockedToday \(appState.stats.blockedToday)"))
                            Text(String(localized: "dashboard.blockedTotal \(appState.stats.blockedTotal)"))
                        }
                        .frame(maxWidth: .infinity, alignment: .leading)
                    }
                    if !appState.stats.byCategory.isEmpty {
                        Chart {
                            ForEach(appState.stats.byCategory.sorted(by: { $0.key < $1.key }), id: \.key) { key, value in
                                BarMark(x: .value("Category", key), y: .value("Count", value))
                            }
                        }
                        .frame(height: 200)
                        .padding()
                    }
                }
                .padding()
            }
            .navigationTitle(String(localized: "tab.stats"))
            .task { await appState.refreshStats() }
        }
    }
}
