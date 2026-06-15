import SwiftUI
import SharedModels
#if canImport(UIKit)
import UIKit
#endif
#if canImport(AppKit)
import AppKit
#endif

public enum MGTheme {
    public static let accent = Color.accentColor

    public static var background: Color {
        #if os(iOS)
        Color(uiColor: .systemGroupedBackground)
        #else
        Color(nsColor: .windowBackgroundColor)
        #endif
    }

    public static var cardBackground: Color {
        #if os(iOS)
        Color(uiColor: .secondarySystemGroupedBackground)
        #else
        Color(nsColor: .controlBackgroundColor)
        #endif
    }
}

public struct MGCard<Content: View>: View {
    private let content: Content

    public init(@ViewBuilder content: () -> Content) {
        self.content = content()
    }

    public var body: some View {
        content
            .padding()
            .background(MGTheme.cardBackground)
            .clipShape(RoundedRectangle(cornerRadius: 12))
    }
}

public struct MGPrimaryButton: View {
    let title: String
    let action: () -> Void
    @Environment(\.userMode) private var userMode

    public init(_ title: String, action: @escaping () -> Void) {
        self.title = title
        self.action = action
    }

    public var body: some View {
        Button(action: action) {
            Text(LocalizedStringKey(title))
                .font(userMode == .elder ? .title3.bold() : .body.bold())
                .frame(maxWidth: .infinity)
                .padding()
                .background(Color.accentColor)
                .foregroundStyle(.white)
                .clipShape(RoundedRectangle(cornerRadius: 12))
        }
    }
}

private struct UserModeKey: EnvironmentKey {
    static let defaultValue: UserMode = .standard
}

public extension EnvironmentValues {
    var userMode: UserMode {
        get { self[UserModeKey.self] }
        set { self[UserModeKey.self] = newValue }
    }
}

public struct UserModeModifier: ViewModifier {
    let mode: UserMode

    public func body(content: Content) -> some View {
        content
            .environment(\.userMode, mode)
            .dynamicTypeSize(mode == .elder ? .xxxLarge ... .accessibility3 : .medium ... .xxxLarge)
    }
}

public extension View {
    func userMode(_ mode: UserMode) -> some View {
        modifier(UserModeModifier(mode: mode))
    }
}
