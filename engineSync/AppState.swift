import Foundation

class AppState: ObservableObject {
    @Published var engineBackupPath: String = ""
    @Published var soundSwitchProjectPath: String = ""
    @Published var targetDrivePath: String = ""

    @Published var loading: Bool = false
    @Published var statusText: String = "Select Engine Backup folder to start!"
    @Published var progressTotalCount: Double = -1
    @Published var progressCurrentIdx: Double = -1
}
