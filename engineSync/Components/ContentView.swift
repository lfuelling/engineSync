import SwiftUI

struct ContentView: View {
    @EnvironmentObject var appState: AppState

    var body: some View {
        VStack {
            Text(appState.statusText)
            ButtonArea().environmentObject(appState)
            Divider()
            SyncButton().environmentObject(appState)
        }.padding()
    }
}

struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        ContentView()
    }
}
