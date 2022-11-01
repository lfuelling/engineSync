package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	fdia "fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
	"image/color"
	"strconv"
	"strings"
)

var engineDbFiles []string
var soundSwitchProject string
var targetDevicePath string

var libraryButton *widget.Button
var soundSwitchButton *widget.Button
var targetDriveButton *widget.Button
var startSyncButton *widget.Button
var statusText *canvas.Text
var loader *widget.ProgressBarInfinite

func setLoading(loading bool) {
	if loading {
		loader.Show()
		loader.Start()
		statusText.Color = color.White
		statusText.Text = "Loading..."
		statusText.Refresh()
		libraryButton.Disable()
		soundSwitchButton.Disable()
		targetDriveButton.Disable()
		startSyncButton.Disable()
	} else {
		loader.Stop()
		loader.Hide()
		statusText.Color = color.White
		statusText.Text = "Ready!"
		statusText.Refresh()

		if !(engineDbFiles != nil && !(len(engineDbFiles) <= 0)) {
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

func main() {
	a := app.New()
	w := a.NewWindow("engineSync")
	w.Resize(fyne.NewSize(255, 255))

	statusText = canvas.NewText("Select Engine Library to start!", color.White)
	statusText.Alignment = fyne.TextAlignCenter

	loader = widget.NewProgressBarInfinite()
	loader.Resize(fyne.NewSize(255, 16))
	loader.Refresh()
	loader.Stop()
	loader.Hide()

	startSyncButton = widget.NewButton("Start Sync!", func() {
		setLoading(true)
	})
	startSyncButton.Disable()

	targetDriveButton = widget.NewButton("Select Target Drive", func() {
		setLoading(true)

		targetPath, err := dialog.Directory().Title("Target Device").Browse()
		if err != nil {
			fdia.NewError(err, w)
		}

		if IsDirectory(targetPath, w) {
			targetDevicePath = targetPath

			splitPath := strings.Split(targetPath, "/")
			targetDriveButton.SetText(splitPath[len(splitPath)-1])

			setLoading(false)
		} else {
			setLoading(false)
			statusText.Color = color.RGBA{R: 255, G: 0, B: 0, A: 1}
			statusText.Text = "Invalid target!"
			statusText.Refresh()
		}
	})
	targetDriveButton.Disable()

	soundSwitchButton = widget.NewButton("Select SoundSwitch Project (Optional)", func() {

		soundSwitchPath, err := dialog.Directory().Title("SoundSwitch Project").Browse()
		if err != nil {
			fdia.NewError(err, w)
		}

		go func() {
			setLoading(true)
			if strings.HasSuffix(soundSwitchPath, ".ssproj") {
				soundSwitchProject = soundSwitchPath
				splitPath := strings.Split(soundSwitchProject, "/")

				soundSwitchButton.SetText(splitPath[len(splitPath)-1])
				setLoading(false)
			} else {
				setLoading(false)
				statusText.Color = color.RGBA{R: 255, G: 0, B: 0, A: 1}
				statusText.Text = "Invalid SoundSwitch Project!"
				statusText.Refresh()
			}
		}()

	})
	soundSwitchButton.Disable()

	libraryButton = widget.NewButton("Select Engine Library", func() {

		directory, err := dialog.Directory().Title("Engine Library").Browse()
		if err != nil {
			fdia.NewError(err, w)
		}

		go func() {
			setLoading(true)

			files := ReadDir(directory, w)
			for _, file := range files {
				if strings.HasSuffix(file, "m.db") {
					engineDbFiles = append(engineDbFiles, file)
				}
			}

			if len(engineDbFiles) > 1 {
				libraryButton.SetText("Found " + strconv.Itoa(len(engineDbFiles)) + " DB files!")
				setLoading(false)
			} else {
				setLoading(false)
				statusText.Color = color.RGBA{R: 255, G: 0, B: 0, A: 1}
				statusText.Text = "Invalid EngineDJ Library!"
				statusText.Refresh()
			}
		}()

	})
	w.SetContent(container.NewVBox(
		layout.NewSpacer(),
		statusText,
		layout.NewSpacer(),
		libraryButton,
		soundSwitchButton,
		targetDriveButton,
		layout.NewSpacer(),
		startSyncButton,
		layout.NewSpacer(),
		loader,
		layout.NewSpacer(),
	))

	w.ShowAndRun()
}
