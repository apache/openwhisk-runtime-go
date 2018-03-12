import SwiftyJSON

#if os(Linux)
import Glibc
#endif

func hello(json: JSON) -> JSON {
    let name = json["name"].string ?? "stranger"
    return JSON([ "greeting" : "Hello \(name)!" ])
}

func printJson(json: JSON) {
    let res = json.rawString()?.replacingOccurrences(of: "\n", with: " ")
    print( res ?? "{\"error:\": \"Invalid JSON\"}")
#if os(Linux)
    fflush(stdout)
#endif
}

func parseJson(str: String) -> JSON {
    if let data = str.data(using: .utf8, allowLossyConversion: false) {
        return JSON(data: data)
    } else {
        return JSON()
    }
}

if CommandLine.arguments.count >= 2 {
    printJson(json: hello(json: parseJson(str: CommandLine.arguments[1])))
} else {
    printJson(json: ["openwhisk", 1])
    while let input = readLine() {
        printJson(json: hello(json: parseJson(str: input)))
    }
}
