package buttons

import (
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
	"image/color"
	"path/filepath"
	"strings"
)

func OnSoundSwitchButtonClick(setLoading func(loading bool, infinite bool, current int, total int),
	setStatus func(status string, color color.RGBA),
	setSoundSwitchPath func(path string),
	soundSwitchButton *widget.Button) {
	soundSwitchPath, err := dialog.Directory().Title("SoundSwitch Project").Browse()
	if err != nil {
		dialog.Message("%s", err).Title("Error!").Error()
		setLoading(false, true, 0, 0)
		return
	}

	go func() {
		setLoading(true, true, 2, 10)
		if strings.HasSuffix(soundSwitchPath, ".ssproj") {
			setSoundSwitchPath(soundSwitchPath)

			splitPath := strings.Split(soundSwitchPath, string(filepath.Separator))
			soundSwitchButton.SetText(splitPath[len(splitPath)-1])
			setLoading(false, true, 0, 0)
		} else {
			setLoading(false, true, 0, 0)
			setStatus("Invalid SoundSwitch Project!", color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}()
}
