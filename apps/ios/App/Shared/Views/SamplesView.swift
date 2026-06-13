import DesignSystem
import SharedModels
import SwiftUI

struct SamplesView: View {
    @Environment(AppState.self) private var appState
    @State private var sampleText = ""
    @State private var selectedLabel: MessageCategory = .spam

    var body: some View {
        NavigationStack {
            Form {
                Section(String(localized: "samples.feed")) {
                    TextField(String(localized: "samples.paste"), text: $sampleText, axis: .vertical)
                        .lineLimit(3 ... 6)
                    Picker(String(localized: "samples.label"), selection: $selectedLabel) {
                        ForEach(MessageCategory.allCases, id: \.self) { cat in
                            Text(cat.rawValue).tag(cat)
                        }
                    }
                    MGPrimaryButton(String(localized: "samples.submit")) {
                        guard !sampleText.isEmpty else { return }
                        Task {
                            await appState.submitSample(text: sampleText, label: selectedLabel)
                            sampleText = ""
                        }
                    }
                }
                Section(String(localized: "samples.history")) {
                    ForEach(appState.samples) { sample in
                        VStack(alignment: .leading) {
                            Text(sample.text).lineLimit(2)
                            Text(sample.label.rawValue).font(.caption).foregroundStyle(.secondary)
                        }
                    }
                }
            }
            .navigationTitle(String(localized: "tab.samples"))
        }
    }
}
