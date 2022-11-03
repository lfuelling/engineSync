package buttons

import (
	"database/sql"
	"engineSync/internal"
	"fmt"
	cp "github.com/otiai10/copy"
	"github.com/sqweek/dialog"
	"image/color"
	"os"
	"path/filepath"
	"strings"
)

func OnSyncButtonClick(setLoading func(loading bool, infinite bool, current int, total int),
	setProgress func(current int),
	setStatus func(status string, color color.RGBA),
	soundSwitchProject string,
	targetDevicePath string,
	engineLibraryDir string,
	engineDbFiles []string) {
	setLoading(true, true, 1, 10)

	if len(soundSwitchProject) > 1 && strings.HasSuffix(soundSwitchProject, ".ssproj") {
		setStatus("Copying SoundSwitch Project", color.RGBA{R: 255, G: 255, B: 255, A: 255})
		targetPath := fmt.Sprintf("%v%vSoundSwitch", targetDevicePath, string(filepath.Separator))
		err := os.MkdirAll(targetPath, os.ModePerm)
		if err != nil {
			dialog.Message("%s", err).Title("Error!").Error()
			setLoading(false, true, 0, 0)
			return
		} else {
			err1 := cp.Copy(soundSwitchProject, targetPath)
			if err1 != nil {
				dialog.Message("%s", err1).Title("Error!").Error()
				setLoading(false, true, 0, 0)
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
			dialog.Message("%s", err).Title("Error!").Error()
			setLoading(false, true, 0, 0)
			return
		} else {
			err := internal.CopyFile(engineDbFile, targetPath)
			if err != nil {
				dialog.Message("%s", err).Title("Error!").Error()
				setLoading(false, true, 0, 0)
				return
			}

			// read tracks from database (from target)
			dbName := splitPath[len(splitPath)-1]
			if dbName == "m.db" { // m.db seems to contain the relevant entries
				setStatus(fmt.Sprintf("Reading tracks of %v...", dbName), color.RGBA{R: 255, G: 255, B: 255, A: 255})

				// open db
				db, err := sql.Open("sqlite", fmt.Sprintf("file:%v%v%v", targetPath, string(filepath.Separator), dbName))
				if err != nil {
					dialog.Message("%s", err).Title("Error!").Error()
					setLoading(false, true, 0, 0)
					return
				}

				// read tracks
				tracks, err2 := internal.GetTracks(db)
				if err2 != nil {
					dialog.Message("%s", err2).Title("Error!").Error()
					setLoading(false, true, 0, 0)
					return
				}

				total := len(tracks)
				setLoading(true, false, 0, total)
				for idx, track := range tracks {
					setProgress(idx)
					setStatus(fmt.Sprintf("Syncing track %v/%v", idx, total), color.RGBA{R: 255, G: 255, B: 255, A: 255})

					err3 := internal.CopyTrack(track, targetDevicePath, engineLibraryDir)
					if err3 != nil {
						dialog.Message("%s", err3).Title("Error!").Error()
						setLoading(false, true, 0, 0)
						return
					}

					err4 := internal.UpdateTrack(track, db)
					if err4 != nil {
						dialog.Message("%s", err4).Title("Error!").Error()
						setLoading(false, true, 0, 0)
						return
					}
				}

				// close db connection
				err4 := db.Close()
				if err4 != nil {
					dialog.Message("%s", err4).Title("Error!").Error()
					setLoading(false, true, 0, 0)
					return
				}
			}
		}
	}
	// all done
	setLoading(false, true, 0, 0)
	setStatus("Finished!", color.RGBA{R: 0, G: 255, B: 0, A: 255})
}
