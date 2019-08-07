package myfdbstorage

import (
	"os"
)

// FileExists check file on disk
func FileExists(file string) (bool, error) {
	_, err := os.Stat(file)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
