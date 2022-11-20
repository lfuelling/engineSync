import SwiftUI

struct SyncButton: View {
    @EnvironmentObject var appState: AppState

    @State private var showError = false

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
                    Task {
                        appState.loading = true
                        appState.statusText = "Creating target directories..."

                        do {
                            try FileManager.default.createDirectory(atPath: "\(appState.targetDrivePath)/Engine Library/Database2", withIntermediateDirectories: true, attributes: nil)
                            try FileManager.default.createDirectory(atPath: "\(appState.targetDrivePath)/SoundSwitch", withIntermediateDirectories: true, attributes: nil)
                        } catch {
                            showError = true
                            print(error)
                        }

                        appState.statusText = "Finished!"
                        appState.loading = false
                    }
                }, label: { Text("Sync").frame(maxWidth: .infinity) })
                        .disabled(syncButtonDisabled())
                        .alert("An error occurred!", isPresented: $showError) {
                            Button("OK", role: .cancel) {
                                showError = false
                            }
                        }
            }
        }
                .padding()
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
