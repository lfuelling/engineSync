import SwiftUI

@main
struct engineSyncApp: App {
    @StateObject var appState: AppState = AppState()

    var body: some Scene {
        WindowGroup {
            ContentView().environmentObject(appState).frame(minWidth: 400, maxWidth: 400, alignment: .top)
        }.windowResizability(.contentSize)
    }
}
