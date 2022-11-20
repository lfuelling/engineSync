import SwiftUI

struct SyncButton: View {
    @EnvironmentObject var appState: AppState

    var body: some View {
        VStack {
            if appState.loading {
                if appState.progressTotalCount > 0 && appState.progressCurrentIdx > 0 {
                    ProgressView(value: appState.progressCurrentIdx, total: appState.progressTotalCount)
                } else {
                    ProgressView()
                }
            } else {
                Button(action: {
                    appState.loading = true
                }, label: {Text("Sync").frame(maxWidth: .infinity)}).disabled(syncButtonDisabled())
            }
        }.padding()
    }

    func syncButtonDisabled() -> Bool {
        // TODO: add proper validation
        appState.engineBackupPath == "" || appState.targetDrivePath == ""
    }
}

struct SyncButton_Previews: PreviewProvider {
    static var previews: some View {
        SyncButton()
    }
}
