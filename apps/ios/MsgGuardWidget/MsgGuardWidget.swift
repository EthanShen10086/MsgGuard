import SharedModels
import SwiftUI
import WidgetKit

struct BlockedStatsEntry: TimelineEntry {
    let date: Date
    let blockedToday: Int
}

struct BlockedStatsProvider: TimelineProvider {
    func placeholder(in context: Context) -> BlockedStatsEntry {
        BlockedStatsEntry(date: .now, blockedToday: 0)
    }

    func getSnapshot(in context: Context, completion: @escaping (BlockedStatsEntry) -> Void) {
        completion(BlockedStatsEntry(date: .now, blockedToday: readCount()))
    }

    func getTimeline(in context: Context, completion: @escaping (Timeline<BlockedStatsEntry>) -> Void) {
        let entry = BlockedStatsEntry(date: .now, blockedToday: readCount())
        completion(Timeline(entries: [entry], policy: .after(.now.addingTimeInterval(900))))
    }

    private func readCount() -> Int {
        UserDefaults(suiteName: AppConstants.appGroupID)?
            .integer(forKey: AppConstants.UserDefaultsKeys.blockedCountToday) ?? 0
    }
}

struct BlockedStatsWidgetView: View {
    let entry: BlockedStatsEntry

    var body: some View {
        VStack(alignment: .leading) {
            Text("MsgGuard")
                .font(.caption)
                .foregroundStyle(.secondary)
            Text("\(entry.blockedToday)")
                .font(.largeTitle.bold())
            Text(String(localized: "widget.blockedToday"))
                .font(.caption)
        }
        .containerBackground(.fill.tertiary, for: .widget)
    }
}

@main
struct MsgGuardWidget: Widget {
    let kind = "MsgGuardWidget"

    var body: some WidgetConfiguration {
        StaticConfiguration(kind: kind, provider: BlockedStatsProvider()) { entry in
            BlockedStatsWidgetView(entry: entry)
        }
        .configurationDisplayName(String(localized: "widget.title"))
        .description(String(localized: "widget.description"))
        .supportedFamilies([.systemSmall, .accessoryInline])
    }
}
