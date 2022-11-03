package internal

import "fyne.io/fyne/v2/widget"

func ShowLoader(loader *widget.ProgressBarInfinite) {
	loader.Start()
	loader.Show()
	loader.Refresh()
}

func HideLoader(loader *widget.ProgressBarInfinite) {
	loader.Stop()
	loader.Hide()
	loader.Refresh()
}

func ShowProgressBar(progress *widget.ProgressBar, total int, current int) {
	progress.Show()
	progress.Min = 0
	progress.Max = float64(total)
	progress.Value = float64(current)
	progress.Refresh()
}

func HideProgressBar(progress *widget.ProgressBar) {
	progress.Min = 0
	progress.Max = 100
	progress.Value = 0
	progress.Hide()
	progress.Refresh()
}
