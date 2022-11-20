//
// Created by Lukas Fülling on 19/11/2022.
//

import SwiftUI

struct FolderSelectorButton: View {
    var label: Text
    @Binding var selectedFolder: String

    var body: some View {
        VStack {
            if selectedFolder == "" {
                Text("No folder selected").foregroundColor(Color.secondary).frame(maxWidth: .infinity)
            } else {
                Text(URL(string: selectedFolder)!.lastPathComponent)
            }
            Button(action: {
                selectFolder()
            }, label: { label.frame(maxWidth: .infinity) })
        }
    }

    func selectFolder() {
        let folderChooserPoint = CGPoint(x: 0, y: 0)
        let folderChooserSize = CGSize(width: 500, height: 600)
        let folderChooserRectangle = CGRect(origin: folderChooserPoint, size: folderChooserSize)
        let folderPicker = NSOpenPanel(contentRect: folderChooserRectangle, styleMask: .utilityWindow, backing: .buffered, defer: true)

        folderPicker.canChooseDirectories = true
        folderPicker.canChooseFiles = false
        folderPicker.allowsMultipleSelection = false
        folderPicker.canDownloadUbiquitousContents = false
        folderPicker.canResolveUbiquitousConflicts = false

        folderPicker.begin { response in

            if response == .OK {
                selectedFolder = folderPicker.url?.absoluteString ?? ""
            }
        }
    }
}
