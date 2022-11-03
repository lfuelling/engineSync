package buttons

import (
	"engineSync/internal"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
	"image/color"
	"path/filepath"
	"strings"
)

func OnTargetDriveButtonClick(setLoading func(loading bool, infinite bool, current int, total int),
	setStatus func(status string, color color.RGBA),
	setTargetPath func(path string),
	targetDriveButton *widget.Button) {
	setLoading(true, true, 2, 10)

	targetPath, err := dialog.Directory().Title("Target Device").Browse()
	if err != nil {
		dialog.Message("%s", err).Title("Error!").Error()
		setLoading(false, true, 0, 0)
		return
	}

	isDirectory, err := internal.IsDirectory(targetPath)
	if err != nil {
		dialog.Message("%s", err).Title("Error!").Error()
		setLoading(false, true, 0, 0)
		return
	}
	if isDirectory {
		setTargetPath(targetPath)

		splitPath := strings.Split(targetPath, string(filepath.Separator))
		targetDriveButton.SetText(splitPath[len(splitPath)-1])

		setLoading(false, true, 0, 0)
	} else {
		setLoading(false, true, 0, 0)
		setStatus("Invalid target!", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	}
}
