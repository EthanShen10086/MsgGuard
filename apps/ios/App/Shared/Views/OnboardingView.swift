import DesignSystem
import SwiftUI

struct OnboardingView: View {
    @Environment(AppState.self) private var appState
    @State private var step = 0

    var body: some View {
        VStack(spacing: 24) {
            Spacer()
            Image(systemName: steps[step].icon)
                .font(.system(size: 64))
                .foregroundStyle(Color.accentColor)
            Text(steps[step].title)
                .font(.title.bold())
                .multilineTextAlignment(.center)
            Text(steps[step].subtitle)
                .font(.body)
                .foregroundStyle(.secondary)
                .multilineTextAlignment(.center)
                .padding(.horizontal)
            Spacer()
            if step == steps.count - 1 {
                MGPrimaryButton(String(localized: "onboarding.openSettings")) {
                    if let url = URL(string: UIApplication.openSettingsURLString) {
                        UIApplication.shared.open(url)
                    }
                }
                Toggle(String(localized: "onboarding.confirmEnabled"), isOn: Binding(
                    get: { appState.extensionEnabled },
                    set: { appState.markExtensionEnabled($0) }
                ))
                .padding(.horizontal)
            }
            MGPrimaryButton(step < steps.count - 1 ? "onboarding.next" : "onboarding.finish") {
                if step < steps.count - 1 {
                    step += 1
                } else {
                    appState.completeOnboarding()
                }
            }
            .padding(.horizontal)
            .padding(.bottom)
        }
    }

    private var steps: [(icon: String, title: String, subtitle: String)] {
        [
            ("shield.lefthalf.filled", String(localized: "onboarding.step1.title"), String(localized: "onboarding.step1.subtitle")),
            ("message.badge.fill", String(localized: "onboarding.step2.title"), String(localized: "onboarding.step2.subtitle")),
            ("gearshape.2.fill", String(localized: "onboarding.step3.title"), String(localized: "onboarding.step3.subtitle")),
        ]
    }
}
