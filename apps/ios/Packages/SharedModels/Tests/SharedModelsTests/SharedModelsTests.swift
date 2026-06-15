import XCTest
@testable import SharedModels

final class SharedModelsTests: XCTestCase {
    func testAppConstantsBundleID() {
        XCTAssertEqual(AppConstants.bundleID, "com.ethanshen.msgguard")
        XCTAssertEqual(AppConstants.appGroupID, "group.com.ethanshen.msgguard")
    }

    func testFilterConfigDefaults() {
        let config = FilterConfig()
        XCTAssertEqual(config.locale, "zh-Hans")
        XCTAssertTrue(config.otpProtectionEnabled)
    }
}
