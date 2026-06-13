// swift-tools-version: 6.0

import PackageDescription

let package = Package(
    name: "FilterEngine",
    platforms: [.iOS(.v17), .macOS(.v14)],
    products: [
        .library(name: "FilterEngine", targets: ["FilterEngine"]),
    ],
    dependencies: [
        .package(path: "../SharedModels"),
    ],
    targets: [
        .target(
            name: "FilterEngine",
            dependencies: ["SharedModels"]
        ),
        .testTarget(
            name: "FilterEngineTests",
            dependencies: ["FilterEngine"]
        ),
    ]
)
