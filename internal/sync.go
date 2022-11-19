package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"
)

func CreateTargetDirectories(targetPath string,
	setLoading func(loading bool, infinite bool, current int, total int),
	setStatus func(status string, color color.RGBA)) (string, string, error) {

	setLoading(true, true, 0, 1)
	setStatus("Creating target directories", color.RGBA{R: 255, G: 255, B: 255, A: 255})

	engineDataPath := fmt.Sprintf("%v%vEngine Library%vDatabase2", targetPath, string(filepath.Separator), string(filepath.Separator))
	soundSwitchDataPath := fmt.Sprintf("%v%vSoundSwitch", targetPath, string(filepath.Separator))

	targetPaths := []string{
		engineDataPath,
		soundSwitchDataPath,
	}

	for _, path := range targetPaths {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return "", "", err
		}
	}

	return engineDataPath, soundSwitchDataPath, nil
}

func CopyEngineDbFiles(setLoading func(loading bool, infinite bool, current int, total int), setProgress func(current int), setStatus func(status string, color color.RGBA), engineDbFiles []string, engineLibraryDir string, engineDataPath string, targetDevicePath string, ignoreNonExistentTracks bool, keepDirectoryStructure bool) (error, []Track) {
	for _, engineDbFile := range engineDbFiles {
		splitPath := strings.Split(engineDbFile, string(filepath.Separator))
		setStatus(fmt.Sprintf("Copying db %v...", splitPath[len(splitPath)-1]), color.RGBA{R: 255, G: 255, B: 255, A: 255})

		// copy db to target
		err := CopyFile(fmt.Sprintf("%v%v%v", engineLibraryDir, string(filepath.Separator), "Database2"), engineDbFile, engineDataPath)
		if err != nil {
			return err, []Track{}
		}

		dbName := splitPath[len(splitPath)-1]
		if dbName == "m.db" { // m.db seems to contain the relevant entries
			setStatus(fmt.Sprintf("Reading tracks of %v...", dbName), color.RGBA{R: 255, G: 255, B: 255, A: 255})

			// open db
			dbPath := fmt.Sprintf("file:%v%v%v", engineDataPath, string(filepath.Separator), dbName)
			db, err := sql.Open("sqlite", dbPath)
			if err != nil {
				return err, []Track{}
			}

			// read tracks
			tracks, err2 := GetTracks(db)
			if err2 != nil {
				return err2, []Track{}
			}

			total := len(tracks)
			setLoading(true, false, 0, total)
			var skippedTracks []Track
			for idx, track := range tracks {
				setProgress(idx)
				setStatus(fmt.Sprintf("Syncing track %v/%v", idx, total), color.RGBA{R: 255, G: 255, B: 255, A: 255})

				// copy tracks
				err3 := CopyTrack(track, targetDevicePath, engineLibraryDir, keepDirectoryStructure)
				if err3 != nil {
					if ignoreNonExistentTracks && errors.Is(err3, os.ErrNotExist) {
						skippedTracks = append(skippedTracks, track)
						continue
					} else {
						return err3, []Track{}
					}
				}

				// Update track in target db
				err4 := UpdateTrack(track, db, keepDirectoryStructure)
				if err4 != nil {
					return err4, []Track{}
				}
			}

			// close db connection
			err4 := db.Close()
			if err4 != nil {
				return err4, []Track{}
			}

			return nil, skippedTracks
		}
	}
	return nil, []Track{}
}
