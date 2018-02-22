// swift-tools-version:4.0
// The swift-tools-version declares the minimum version of Swift required to build this package.

import PackageDescription

let package = Package(
    name: "Hello",
    products:[
        .executable(name: "Hello", targets: ["Hello"]),
        .library(name: "Hello", type: .dynamic, targets: ["Hello"])
    ],
    dependencies: [
        .package(url: "https://github.com/IBM-Swift/SwiftyJSON", from: "17.0.0"),
    ],
    targets: [
        // Targets are the basic building blocks of a package. A target can define a module or a test suite.
        // Targets can depend on other targets in this package, and on products in packages which this package depends on.
        .target(
            name: "Hello",
            dependencies: ["SwiftyJSON"]),
    ]
)
