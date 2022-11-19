package main

import (
	"engineSync/buttons"
	"engineSync/internal"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	_ "modernc.org/sqlite"
	"strings"
)

var engineLibraryDir string
var engineDbFiles []string
var soundSwitchProject string
var targetDevicePath string
var ignoreNonExistentTracks bool
var keepDirectoryStructure bool

var libraryButton *widget.Button
var soundSwitchButton *widget.Button
var targetDriveButton *widget.Button
var startSyncButton *widget.Button
var statusText *canvas.Text
var ignoreNonExistentTracksCheck *widget.Check
var keepDirectoryStructureCheck *widget.Check
var progress *widget.ProgressBar
var loader *widget.ProgressBarInfinite

func main() {
	a := app.New()
	w := a.NewWindow("engineSync")
	w.Resize(fyne.NewSize(255, 280))
	w.SetFixedSize(true)

	statusText = canvas.NewText("Select Engine Library to start!", color.White)
	statusText.Alignment = fyne.TextAlignCenter

	progress = widget.NewProgressBar()
	progress.Hide()
	progress.Refresh()

	loader = widget.NewProgressBarInfinite()
	loader.Hide()
	loader.Refresh()

	startSyncButton = widget.NewButton("Start Sync!", func() {
		buttons.OnSyncButtonClick(setLoading, func(current int) {
			progress.Value = float64(current)
			progress.Refresh()
		}, setStatus, soundSwitchProject, targetDevicePath, engineLibraryDir, engineDbFiles, ignoreNonExistentTracks, keepDirectoryStructure)
	})
	startSyncButton.Disable()

	targetDriveButton = widget.NewButton("Select Target Drive", func() {
		buttons.OnTargetDriveButtonClick(setLoading, setStatus, func(path string) {
			targetDevicePath = path
		}, targetDriveButton)
	})
	targetDriveButton.Disable()

	soundSwitchButton = widget.NewButton("Select SoundSwitch Project (Optional)", func() {
		buttons.OnSoundSwitchButtonClick(setLoading, setStatus, func(path string) {
			soundSwitchProject = path
		}, soundSwitchButton)
	})
	soundSwitchButton.Disable()

	libraryButton = widget.NewButton("Select Engine Library Backup", func() {
		buttons.OnLibraryButtonClick(setLoading, setStatus, func(files []string) {
			engineDbFiles = files
		}, func(path string) {
			engineLibraryDir = path
		}, libraryButton)
	})

	ignoreNonExistentTracksCheck = widget.NewCheck("Ignore Missing Tracks", func(value bool) {
		ignoreNonExistentTracks = value
	})

	keepDirectoryStructureCheck = widget.NewCheck("Keep directory structure", func(value bool) {
		keepDirectoryStructure = value
	})

	w.SetContent(container.NewVBox(
		layout.NewSpacer(),
		statusText,
		layout.NewSpacer(),
		libraryButton,
		soundSwitchButton,
		targetDriveButton,
		container.New(layout.NewCenterLayout(), container.NewVBox(
			ignoreNonExistentTracksCheck,
			keepDirectoryStructureCheck,
		)),
		layout.NewSpacer(),
		startSyncButton,
		layout.NewSpacer(),
		container.NewVBox(
			progress,
			loader,
		),
		layout.NewSpacer(),
	))

	w.ShowAndRun()
}

func setStatus(status string, color color.RGBA) {
	statusText.Color = color
	statusText.Text = status
	statusText.Refresh()
}

func setLoading(loading bool, infinite bool, current int, total int) {
	if loading {
		if infinite {
			internal.HideProgressBar(progress)
			internal.ShowLoader(loader)
		} else {
			internal.HideLoader(loader)
			internal.ShowProgressBar(progress, total, current)
		}

		setStatus("Loading...", color.RGBA{R: 255, G: 255, B: 255, A: 255})
		libraryButton.Disable()
		soundSwitchButton.Disable()
		targetDriveButton.Disable()
		ignoreNonExistentTracksCheck.Disable()
		keepDirectoryStructureCheck.Disable()
		startSyncButton.Disable()
	} else {
		internal.HideProgressBar(progress)
		internal.HideLoader(loader)

		setStatus("Ready!", color.RGBA{R: 255, G: 255, B: 255, A: 255})

		ignoreNonExistentTracksCheck.Enable()
		keepDirectoryStructureCheck.Enable()
		if !(engineDbFiles != nil && !(len(engineDbFiles) <= 0) || len(engineLibraryDir) <= 0) {
			libraryButton.Enable()
		} else {
			if len(soundSwitchProject) <= 0 && !strings.HasSuffix(soundSwitchProject, ".ssproj") {
				soundSwitchButton.Enable()
			}
			if len(targetDevicePath) <= 0 {
				targetDriveButton.Enable()
			} else {
				startSyncButton.Enable()
			}
		}
	}
}
