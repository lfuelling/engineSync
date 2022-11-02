package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	fdia "fyne.io/fyne/v2/dialog"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ReadDir(path string, w fyne.Window) []string {
	var result []string
	file, err := os.Open(path)
	if err != nil {
		fdia.NewError(err, w).Show()
	}
	defer file.Close()
	names, _ := file.Readdirnames(0)
	for _, name := range names {
		filePath := fmt.Sprintf("%v%v%v", path, string(filepath.Separator), name)
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
		fdia.NewError(err, w).Show()
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		fdia.NewError(err, w).Show()
	}

	return fileInfo.IsDir()
}

func CopyFile(src string, dstDir string, w fyne.Window) int64 {
	fin, err := os.Open(src)
	if err != nil {
		fdia.NewError(err, w).Show()
	}
	defer fin.Close()

	splitPath := strings.Split(src, string(filepath.Separator))
	fout, err2 := os.Create(fmt.Sprintf("%v%v%v", dstDir, string(filepath.Separator), splitPath[len(splitPath)-1]))
	if err2 != nil {
		fdia.NewError(err2, w).Show()
	}
	defer fout.Close()

	written, err3 := io.Copy(fout, fin)

	if err3 != nil {
		fdia.NewError(err3, w).Show()
	}

	return written
}

func CopyTrack(src string, dstDir string) (string, error) {
	fin, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer fin.Close()

	splitPath := strings.Split(src, string(filepath.Separator))
	fileName := splitPath[len(splitPath)-1]
	fout, err2 := os.Create(fmt.Sprintf("%v%v%v", dstDir, string(filepath.Separator), fileName))
	if err2 != nil {
		return "", err2
	}
	defer fout.Close()

	_, err3 := io.Copy(fout, fin)
	if err3 != nil {
		return "", err3
	}

	return fileName, nil
}
