package main

import (
	"database/sql"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	fdia "fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	cp "github.com/otiai10/copy"
	"github.com/sqweek/dialog"
	"image/color"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
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

func setStatus(status string, color color.RGBA) {
	statusText.Color = color
	statusText.Text = status
	statusText.Refresh()
}

func setLoading(loading bool) {
	if loading {
		loader.Show()
		loader.Start()
		setStatus("Loading...", color.RGBA{R: 255, G: 255, B: 255, A: 255})
		libraryButton.Disable()
		soundSwitchButton.Disable()
		targetDriveButton.Disable()
		startSyncButton.Disable()
	} else {
		loader.Stop()
		loader.Hide()
		setStatus("Ready!", color.RGBA{R: 255, G: 255, B: 255, A: 255})

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

func handleTrack(rows *sql.Rows, db *sql.DB) error {
	track := Track{}
	err2 := rows.Scan(&track.Id, &track.Path, &track.Filename)
	if err2 != nil {
		return err2
	}

	// copy track to target Music folder
	splitTrackPath := strings.Split(track.Path, "#")[0]
	splitTrackPathSegments := strings.Split(splitTrackPath, string(filepath.Separator))
	setStatus(fmt.Sprintf("Syncing track %v...", splitTrackPathSegments[len(splitTrackPathSegments)-1]), color.RGBA{R: 255, G: 255, B: 255, A: 255})

	targetPath := fmt.Sprintf("%v%vEngine Library%vMusic", targetDevicePath, string(filepath.Separator), string(filepath.Separator))
	err3 := os.MkdirAll(targetPath, os.ModePerm)
	if err3 != nil {
		return err3
	}
	trackPath, err4 := filepath.Abs(splitTrackPath)
	if err4 != nil {
		return err4
	}

	_, err5 := CopyTrack(trackPath, targetPath)
	if err5 != nil {
		return err5
	}
	split := strings.Split(track.Path, string(filepath.Separator)) //FIXME: the path seems to contain "garbage"
	newPath := fmt.Sprintf("Music%v%v", string(filepath.Separator), split[len(split)-1])
	_, err6 := db.Exec("UPDATE Track SET path = ? WHERE id = ?;", newPath, track.Id)
	if err6 != nil {
		return err6
	}
	return nil
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

		if len(soundSwitchProject) > 1 && strings.HasSuffix(soundSwitchProject, ".ssproj") {
			setStatus("Copying SoundSwitch Project", color.RGBA{R: 255, G: 255, B: 255, A: 255})
			targetPath := fmt.Sprintf("%v%vSoundSwitch", targetDevicePath, string(filepath.Separator))
			err := os.MkdirAll(targetPath, os.ModePerm)
			if err != nil {
				fdia.NewError(err, w).Show()
				setLoading(false)
				return
			} else {
				err1 := cp.Copy(soundSwitchProject, targetPath)
				if err1 != nil {
					fdia.NewError(err1, w).Show()
					setLoading(false)
					return
				}
			}
		}

		for _, engineDbFile := range engineDbFiles {
			splitPath := strings.Split(engineDbFile, string(filepath.Separator))
			setStatus(fmt.Sprintf("Copying db %v...", splitPath[len(splitPath)-1]), color.RGBA{R: 255, G: 255, B: 255, A: 255})

			// copy db to target
			targetPath := fmt.Sprintf("%v%vEngine Library%vDatabase2", targetDevicePath, string(filepath.Separator), string(filepath.Separator))
			err := os.MkdirAll(targetPath, os.ModePerm)
			if err != nil {
				fdia.NewError(err, w).Show()
				setLoading(false)
				return
			} else {
				CopyFile(engineDbFile, targetPath, w)

				// read tracks from database (from target)
				dbName := splitPath[len(splitPath)-1]
				setStatus(fmt.Sprintf("Reading tracks of %v...", dbName), color.RGBA{R: 255, G: 255, B: 255, A: 255})
				db, err := sql.Open("sqlite", fmt.Sprintf("file:%v%v%v", targetPath, string(filepath.Separator), dbName))
				if err != nil {
					fdia.NewError(err, w).Show()
					setLoading(false)
					return
				}

				rows, err1 := db.Query("SELECT id, path, filename FROM Track;")
				if err1 != nil {
					fdia.NewError(err1, w).Show()
					setLoading(false)
					return
				}

				for rows.Next() {
					err2 := handleTrack(rows, db)
					if err2 != nil {
						fdia.NewError(err2, w).Show()
						setLoading(false)
						return
					}
				}

				err3 := rows.Close()
				if err3 != nil {
					fdia.NewError(err3, w).Show()
					setLoading(false)
					return
				}

				err4 := db.Close()
				if err4 != nil {
					fdia.NewError(err4, w).Show()
					setLoading(false)
					return
				}
			}
		}
	})
	startSyncButton.Disable()

	targetDriveButton = widget.NewButton("Select Target Drive", func() {
		setLoading(true)

		targetPath, err := dialog.Directory().Title("Target Device").Browse()
		if err != nil {
			fdia.NewError(err, w).Show()
			setLoading(false)
			return
		}

		if IsDirectory(targetPath, w) {
			targetDevicePath = targetPath

			splitPath := strings.Split(targetPath, string(filepath.Separator))
			targetDriveButton.SetText(splitPath[len(splitPath)-1])

			setLoading(false)
		} else {
			setLoading(false)
			setStatus("Invalid target!", color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	})
	targetDriveButton.Disable()

	soundSwitchButton = widget.NewButton("Select SoundSwitch Project (Optional)", func() {

		soundSwitchPath, err := dialog.Directory().Title("SoundSwitch Project").Browse()
		if err != nil {
			fdia.NewError(err, w).Show()
			setLoading(false)
			return
		}

		go func() {
			setLoading(true)
			if strings.HasSuffix(soundSwitchPath, ".ssproj") {
				soundSwitchProject = soundSwitchPath
				splitPath := strings.Split(soundSwitchProject, string(filepath.Separator))

				soundSwitchButton.SetText(splitPath[len(splitPath)-1])
				setLoading(false)
			} else {
				setLoading(false)
				setStatus("Invalid SoundSwitch Project!", color.RGBA{R: 255, G: 0, B: 0, A: 255})
			}
		}()

	})
	soundSwitchButton.Disable()

	libraryButton = widget.NewButton("Select Engine Library", func() {

		directory, err := dialog.Directory().Title("Engine Library").Browse()
		if err != nil {
			fdia.NewError(err, w).Show()
			setLoading(false)
			return
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
				setStatus("Invalid EngineDJ Library!", color.RGBA{R: 255, G: 0, B: 0, A: 255})
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
