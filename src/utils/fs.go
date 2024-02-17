package utils

import (
	"os"
)

func MakeSureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}
func FsIsExist(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		return false
	}
	return true
}