// swift-tools-version: 6.0

import PackageDescription

let package = Package(
    name: "BlocklistStore",
    platforms: [.iOS(.v17), .macOS(.v14)],
    products: [
        .library(name: "BlocklistStore", targets: ["BlocklistStore"]),
    ],
    dependencies: [
        .package(path: "../SharedModels"),
    ],
    targets: [
        .target(
            name: "BlocklistStore",
            dependencies: ["SharedModels"]
        ),
    ]
)
