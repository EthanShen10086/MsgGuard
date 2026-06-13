import SwiftUI

struct HelpView: View {
    var body: some View {
        NavigationStack {
            List {
                Section(String(localized: "help.intro")) {
                    Text(String(localized: "help.intro.body"))
                }
                Section(String(localized: "help.faq")) {
                    Text(String(localized: "help.faq.q1"))
                    Text(String(localized: "help.faq.a1")).foregroundStyle(.secondary)
                    Text(String(localized: "help.faq.q2"))
                    Text(String(localized: "help.faq.a2")).foregroundStyle(.secondary)
                }
                Section(String(localized: "help.feedback")) {
                    NavigationLink(String(localized: "help.feedback.submit")) {
                        FeedbackView()
                    }
                    Link(String(localized: "settings.support"), destination: URL(string: "https://msgguard.app/support")!)
                }
            }
            .navigationTitle(String(localized: "tab.help"))
        }
    }
}
