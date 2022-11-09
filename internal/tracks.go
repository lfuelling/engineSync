package internal

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CopyTrack(track Track, targetDevicePath string, engineLibraryDir string) error {
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
	trackBaseDir, err := getTrackBaseDirAbsolute(track, engineLibraryDir)
	if err != nil {
		return err
	}

	// copy track to target path
	err = CopyFile(trackBaseDir, trackPath, targetPath)
	if err != nil {
		return err
	}

	return nil
}

func getTrackBaseDirRelative(track Track) string {
	splitBySeparator := strings.Split(track.Path, string(filepath.Separator))
	trackDir := strings.Join(splitBySeparator[:len(splitBySeparator)-1], string(filepath.Separator))
	splitTrackDir := strings.Split(trackDir, "..") //TODO: use constant similar to filepath.Separator, not sure if it exists
	splitTrackDirBySeparator := strings.Split(splitTrackDir[len(splitTrackDir)-1], string(filepath.Separator))
	var relativeBaseDirBuilder strings.Builder
	for i, dir := range splitTrackDirBySeparator {
		if dir != "" {
			for j := 0; j < (i/2)+1; j++ {
				relativeBaseDirBuilder.WriteString("..") //TODO: use constant similar to filepath.Separator, not sure if it exists
				relativeBaseDirBuilder.WriteString(string(filepath.Separator))
			}
			relativeBaseDirBuilder.WriteString(dir)
			return relativeBaseDirBuilder.String()
		}
	}

	return ""
}

func getTrackBaseDirAbsolute(track Track, engineLibraryDir string) (string, error) {
	return filepath.Abs(fmt.Sprintf("%v%v%v", engineLibraryDir, string(filepath.Separator), getTrackBaseDirRelative(track)))
}

func UpdateTrack(track Track, db *sql.DB) error {
	// begin new db transaction
	begin, err := db.Begin()
	if err != nil {
		return err
	}

	// update track path in db
	_, err1 := begin.Exec("UPDATE OR REPLACE Track SET path = ? WHERE id = ?;", buildTrackDbPath(track), track.Id)
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

func buildTrackDbPath(track Track) string {
	split := strings.Split(track.Path, string(filepath.Separator))
	fileName := split[len(split)-1]
	relativeBaseDir := getTrackBaseDirRelative(track)
	trackRelativePath := strings.Split(track.Path, fileName)[0]
	trackSubDir := strings.Split(trackRelativePath, relativeBaseDir)[1]

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
