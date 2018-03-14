// swift-tools-version:4.0
// The swift-tools-version declares the minimum version of Swift required to build this package.

import PackageDescription

let package = Package(
    name: "exec",
    products:[
        .executable(name: "exec", targets: ["exec"]),
        .library(name: "exec", type: .dynamic, targets: ["exec"])
    ],
    dependencies: [
        .package(url: "https://github.com/IBM-Swift/SwiftyJSON", from: "17.0.0"),
        .package(url: "https://github.com/IBM-Swift/HeliumLogger", from: "1.7.1"),
       
    ],
    targets: [
        // Targets are the basic building blocks of a package. A target can define a module or a test suite.
        // Targets can depend on other targets in this package, and on products in packages which this package depends on.
        .target(
            name: "exec",
            dependencies: ["SwiftyJSON", "HeliumLogger"]),
    ]
)
