import SwiftUI

struct ButtonArea: View {
    @EnvironmentObject var appState: AppState

    var body: some View {
        VStack(content: {
            FolderSelectorButton(label: Text("Select Engine Backup Folder"), selectedFolder: $appState.engineBackupPath)
            FolderSelectorButton(label: Text("Select SoundSwitch Project (Optional)"), selectedFolder: $appState.soundSwitchProjectPath)
            FolderSelectorButton(label: Text("Select Target Drive"), selectedFolder: $appState.targetDrivePath)
        }).padding()
    }
}

struct ButtonArea_Previews: PreviewProvider {
    static var previews: some View {
        ButtonArea()
    }
}
