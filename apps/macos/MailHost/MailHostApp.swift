import SwiftUI

@main
struct MsgGuardMailHostApp: App {
    var body: some Scene {
        WindowGroup {
            MailHostView()
        }
    }
}

struct MailHostView: View {
    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("MsgGuard Mail")
                .font(.title)
            Text("Enable the Mail extension in Mail → Settings → Extensions, then MsgGuard Mail will classify incoming messages.")
                .foregroundStyle(.secondary)
            Text("Shared App Group: group.com.ethanshen.msgguard")
                .font(.caption)
        }
        .padding(24)
        .frame(minWidth: 420, minHeight: 200)
    }
}
