package internal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ReadDir(path string) ([]string, error) {
	var result []string
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	names, _ := file.Readdirnames(0)
	for _, name := range names {
		filePath := fmt.Sprintf("%v%v%v", path, string(filepath.Separator), name)
		isDirectory, err := IsDirectory(filePath)
		if err != nil {
			return nil, err
		}
		if isDirectory {
			dirResult, err1 := ReadDir(filePath)
			if err1 != nil {
				return nil, err1
			}
			result = append(result, dirResult...)
		} else {
			result = append(result, filePath)
		}
	}
	return result, nil
}

func IsDirectory(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}

func CopyFile(src string, dstDir string) error {
	fin, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fin.Close()

	splitPath := strings.Split(src, string(filepath.Separator))
	fout, err2 := os.Create(fmt.Sprintf("%v%v%v", dstDir, string(filepath.Separator), splitPath[len(splitPath)-1]))
	if err2 != nil {
		return err2
	}
	defer fout.Close()

	_, err3 := io.Copy(fout, fin)

	if err3 != nil {
		return err3
	}

	return nil
}

func copyTrackFile(src string, dstDir string) error {
	fin, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fin.Close()

	splitPath := strings.Split(src, string(filepath.Separator))
	fileName := splitPath[len(splitPath)-1]
	fout, err := os.Create(fmt.Sprintf("%v%v%v", dstDir, string(filepath.Separator), fileName))
	if err != nil {
		return err
	}
	defer fout.Close()

	_, err = io.Copy(fout, fin)
	if err != nil {
		return err
	}

	return nil
}
