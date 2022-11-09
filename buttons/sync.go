package buttons

import (
	"engineSync/internal"
	"github.com/sqweek/dialog"
	"image/color"
)

func OnSyncButtonClick(setLoading func(loading bool, infinite bool, current int, total int),
	setProgress func(current int),
	setStatus func(status string, color color.RGBA),
	soundSwitchProject string,
	targetDevicePath string,
	engineLibraryDir string,
	engineDbFiles []string) {
	setLoading(true, true, 1, 10)

	engineDataPath, soundSwitchDataPath, err := internal.CreateTargetDirectories(targetDevicePath, setLoading, setStatus)
	if err != nil {
		dialog.Message("%s", err).Title("Error!").Error()
		setLoading(false, true, 0, 0)
		return
	}

	if len(soundSwitchProject) > 1 && internal.IsSoundSwitchProjectDir(soundSwitchProject) {
		err1 := internal.SyncSoundSwitchProject(soundSwitchProject, soundSwitchDataPath, setStatus, setProgress, setLoading)
		if err1 != nil {
			dialog.Message("%s", err1).Title("Error!").Error()
			setLoading(false, true, 0, 0)
			return
		}
	}

	err2 := internal.CopyEngineDbFiles(setLoading, setProgress, setStatus, engineDbFiles, engineLibraryDir, engineDataPath, targetDevicePath)
	if err2 != nil {
		dialog.Message("%s", err2).Title("Error!").Error()
		setLoading(false, true, 0, 0)
		return
	}

	// all done
	setLoading(false, true, 0, 0)
	setStatus("Finished!", color.RGBA{R: 0, G: 255, B: 0, A: 255})
}
