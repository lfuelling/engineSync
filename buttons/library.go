package buttons

import (
	"engineSync/internal"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
	"image/color"
	"strconv"
	"strings"
)

func OnLibraryButtonClick(setLoading func(loading bool, infinite bool, current int, total int),
	setStatus func(status string, color color.RGBA),
	setEngineDbFiles func(files []string),
	setEngineLibraryDir func(path string),
	libraryButton *widget.Button) {
	directory, err := dialog.Directory().Title("Engine Library Backup").Browse()
	if err != nil {
		dialog.Message("%s", err).Title("Error!").Error()
		setLoading(false, true, 0, 0)
		return
	}

	go func() {
		setLoading(true, true, 2, 10)

		var engineDbFiles []string

		files, err := internal.ReadDir(directory)
		if err != nil {
			dialog.Message("%s", err).Title("Error!").Error()
			setLoading(false, true, 0, 0)
			return
		}
		for _, file := range files {
			if strings.HasSuffix(file, "m.db") {
				engineDbFiles = append(engineDbFiles, file)
			}
		}

		if len(engineDbFiles) > 1 {
			setEngineDbFiles(engineDbFiles)
			setEngineLibraryDir(directory)
			libraryButton.SetText("Found " + strconv.Itoa(len(engineDbFiles)) + " DB files!")
			setLoading(false, true, 0, 0)
		} else {
			setLoading(false, true, 0, 0)
			setStatus("Invalid EngineDJ Library!", color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}()
}
