package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	fdia "fyne.io/fyne/v2/dialog"
	"os"
)

func ReadDir(path string, w fyne.Window) []string {
	var result []string
	file, err := os.Open(path)
	if err != nil {
		fdia.NewError(err, w)
	}
	defer file.Close()
	names, _ := file.Readdirnames(0)
	for _, name := range names {
		filePath := fmt.Sprintf("%v/%v", path, name)
		if IsDirectory(filePath, w) {
			dirResult := ReadDir(filePath, w)
			result = append(result, dirResult...)
		} else {
			result = append(result, filePath)
		}
	}
	return result
}

func IsDirectory(path string, w fyne.Window) bool {
	file, err := os.Open(path)
	if err != nil {
		fdia.NewError(err, w)
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		fdia.NewError(err, w)
	}

	return fileInfo.IsDir()
}
