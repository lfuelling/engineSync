package internal

import (
	"errors"
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

func CopyFile(baseDir string, src string, dstDir string) error {
	if !strings.HasPrefix(src, baseDir) {
		return errors.New("src must start with baseDir")
	}

	fin, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fin.Close()

	splitBaseDir := strings.Split(src, baseDir)
	splitPath := strings.Split(splitBaseDir[1], string(filepath.Separator))
	fileName := splitPath[len(splitPath)-1]
	targetSubDir := strings.Split(splitBaseDir[1], fileName)[0]
	var destinationPath string
	if targetSubDir == "/" {
		destinationPath = fmt.Sprintf("%v%v%v", dstDir, string(filepath.Separator), fileName)
	} else {
		destinationPathWithSubDir := fmt.Sprintf("%v%v", dstDir, targetSubDir)
		err2 := os.MkdirAll(destinationPathWithSubDir, os.ModePerm)
		if err2 != nil {
			return err2
		}
		destinationPath = fmt.Sprintf("%v%v", destinationPathWithSubDir, fileName)
	}

	fout, err2 := os.Create(destinationPath)
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

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return true, err
}
