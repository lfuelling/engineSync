package internal

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func CopyTrack(track Track, targetDevicePath string, engineLibraryDir string, keepDirectoryStructure bool) error {
	// build target path
	targetPath := fmt.Sprintf("%v%vEngine Library%vMusic", targetDevicePath, string(filepath.Separator), string(filepath.Separator))

	// create target path
	err := os.MkdirAll(targetPath, os.ModePerm)
	if err != nil {
		return err
	}

	// build track path
	trackPath, err := filepath.Abs(fmt.Sprintf("%v%v%v", engineLibraryDir, string(filepath.Separator), track.Path))
	if err != nil {
		return err
	}

	// get track basedir
	var trackBaseDir string
	if keepDirectoryStructure {
		trackBaseDir, err = getTrackBaseDirAbsolute(track, engineLibraryDir)
		if err != nil {
			return err
		}
	} else {
		trackBaseDir, err = filepath.Abs(fmt.Sprintf("%v%v%v", engineLibraryDir, string(filepath.Separator), path.Dir(track.Path)))
		if err != nil {
			return err
		}
	}

	// copy track to target path
	err = CopyFile(trackBaseDir, trackPath, targetPath)
	if err != nil {
		return err
	}

	return nil
}

func getTrackBaseDirRelative(track Track) string {
	var index int
	split := strings.Split(track.Path, "/")
	for i, s := range split {
		if s != ".." {
			index = i
			break
		}
	}
	baseDir := split[:index+1]
	return strings.Join(baseDir, "/")
}

func getTrackBaseDirAbsolute(track Track, engineLibraryDir string) (string, error) {
	return filepath.Abs(fmt.Sprintf("%v%v%v", engineLibraryDir, string(filepath.Separator), getTrackBaseDirRelative(track)))
}

func UpdateTrack(track Track, db *sql.DB, keepDirectoryStructure bool) error {
	// begin new db transaction
	begin, err := db.Begin()
	if err != nil {
		return err
	}

	// update track path in db
	_, err1 := begin.Exec("UPDATE OR REPLACE Track SET path = ? WHERE id = ?;", buildTrackDbPath(track, keepDirectoryStructure), track.Id)
	if err1 != nil {
		err2 := begin.Rollback()
		if err2 != nil {
			return err2
		}
		return err1
	}

	// commit db transaction
	err2 := begin.Commit()
	if err2 != nil {
		return err2
	}

	return nil
}

func buildTrackDbPath(track Track, keepDirectoryStructure bool) string {
	split := strings.Split(track.Path, string(filepath.Separator))
	fileName := split[len(split)-1]

	var trackSubDir string
	if keepDirectoryStructure {
		relativeBaseDir := getTrackBaseDirRelative(track)
		trackRelativePath := strings.Split(track.Path, fileName)[0]
		trackSubDir = strings.Split(trackRelativePath, relativeBaseDir)[1]
	} else {
		trackSubDir = "/"
	}

	newPath := fmt.Sprintf("Music%v%v", trackSubDir, fileName)
	return newPath
}

func GetTracks(db *sql.DB) ([]Track, error) {
	// query db
	rows, err1 := db.Query("SELECT id, path, filename FROM Track;")
	if err1 != nil {
		return nil, err1
	}

	// create result slice
	var tracks []Track

	// iterate rows
	for rows.Next() {
		track := Track{}

		// read db info
		err2 := rows.Scan(&track.Id, &track.Path, &track.Filename)
		if err2 != nil {
			return nil, err2
		}

		// append track
		tracks = append(tracks, track)
	}

	// close db query
	err3 := rows.Close()
	if err3 != nil {
		return nil, err3
	}

	// return result
	return tracks, nil
}
