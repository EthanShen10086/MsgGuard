import DesignSystem
import SharedModels
import SwiftUI

struct RulesView: View {
    @Environment(AppState.self) private var appState
    @State private var newPattern = ""
    @State private var newRuleType: RuleType = .keywordBlock

    var body: some View {
        NavigationStack {
            Form {
                Section(String(localized: "rules.filterTags")) {
                    ForEach(FilterTag.allCases) { tag in
                        Toggle(tag.rawValue, isOn: Binding(
                            get: { appState.filterConfig.enabledTags.contains(tag) },
                            set: { enabled in
                                if enabled {
                                    appState.filterConfig.enabledTags.insert(tag)
                                } else {
                                    appState.filterConfig.enabledTags.remove(tag)
                                }
                            }
                        ))
                    }
                }
                Section(String(localized: "rules.customRules")) {
                    Picker(String(localized: "rules.type"), selection: $newRuleType) {
                        Text("keywordBlock").tag(RuleType.keywordBlock)
                        Text("keywordAllow").tag(RuleType.keywordAllow)
                        Text("numberBlock").tag(RuleType.numberBlock)
                        Text("numberAllow").tag(RuleType.numberAllow)
                    }
                    TextField(String(localized: "rules.pattern"), text: $newPattern)
                    Button(String(localized: "rules.add")) {
                        let rule = BlockRule(type: newRuleType, pattern: newPattern, priority: appState.filterConfig.rules.count)
                        appState.filterConfig.rules.append(rule)
                        newPattern = ""
                    }
                    ForEach(appState.filterConfig.rules) { rule in
                        HStack {
                            Text("\(rule.type.rawValue): \(rule.pattern)")
                            Spacer()
                            Toggle("", isOn: Binding(
                                get: { rule.enabled },
                                set: { val in
                                    if let idx = appState.filterConfig.rules.firstIndex(where: { $0.id == rule.id }) {
                                        appState.filterConfig.rules[idx].enabled = val
                                    }
                                }
                            ))
                            .labelsHidden()
                        }
                    }
                    .onDelete { indices in
                        appState.filterConfig.rules.remove(atOffsets: indices)
                    }
                }
            }
            .navigationTitle(String(localized: "tab.rules"))
            .toolbar {
                ToolbarItem(placement: .topBarTrailing) {
                    Button(String(localized: "rules.save")) {
                        Task { await appState.saveConfig() }
                    }
                }
            }
        }
    }
}
