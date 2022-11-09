package internal

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
)

func IsSoundSwitchProjectDir(path string) bool {
	jsonPath := fmt.Sprintf("%v%v%v", path, string(filepath.Separator), ".ssproj")
	exists, err2 := FileExists(jsonPath)
	return err2 == nil && exists
}

func SyncSoundSwitchProject(
	soundSwitchProject string,
	targetDataPath string,
	setStatus func(status string, color color.RGBA),
	setProgress func(current int),
	setLoading func(loading bool, infinite bool, current int, total int)) error {

	err := copySoundSwitchProjectJSON(soundSwitchProject, targetDataPath, setStatus)
	if err != nil {
		return err
	}

	err2 := copySoundSwitchProjectFiles(soundSwitchProject, targetDataPath, setProgress, setStatus, setLoading)
	return err2
}

func copySoundSwitchProjectFiles(soundSwitchProject string, targetDataPath string, setProgress func(current int), setStatus func(status string, color color.RGBA), setLoading func(loading bool, infinite bool, current int, total int)) error {
	// set status
	setStatus("Copying SoundSwitch files", color.RGBA{R: 255, G: 255, B: 255, A: 255})

	// get list of files
	files, err1 := ReadDir(soundSwitchProject)
	if err1 != nil {
		return err1
	}

	// copy files
	totalFiles := len(files)
	setLoading(true, false, 0, totalFiles)
	for i, file := range files {
		statusString := fmt.Sprintf("Copying SoundSwitch files (%v/%v)", i, totalFiles)
		setStatus(statusString, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		setProgress(i)

		err2 := CopyFile(soundSwitchProject, file, targetDataPath)
		if err2 != nil {
			return err2
		}
	}
	return nil
}

func copySoundSwitchProjectJSON(soundSwitchProject string, targetDataPath string, setStatus func(status string, color color.RGBA)) error {
	// set status
	setStatus("Copying SoundSwitch project", color.RGBA{R: 255, G: 255, B: 255, A: 255})

	// read SoundSwitch project
	jsonBytes, err := os.ReadFile(fmt.Sprintf("%v%v%v", soundSwitchProject, string(filepath.Separator), ".ssproj"))
	if err != nil {
		return err
	}

	// parse into struct
	soundSwitchProjectStruct := SoundSwitchProject{}
	err2 := json.Unmarshal(jsonBytes, &soundSwitchProjectStruct)
	if err2 != nil {
		return err2
	}

	// set read only
	soundSwitchProjectStruct.ReadOnly = true

	// turn it back into []byte
	marshalledData, err3 := json.Marshal(soundSwitchProjectStruct)
	if err3 != nil {
		return err3
	}

	// write to target dir
	err4 := os.WriteFile(fmt.Sprintf("%v%v%v", targetDataPath, string(filepath.Separator), ".ssproj"), marshalledData, os.ModePerm)
	return err4
}
