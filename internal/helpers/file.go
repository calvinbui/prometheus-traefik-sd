package helpers

import (
	"os"
	"path"
)

func InitFolder(filePath string) error {
	folder := path.Dir(filePath)

	_, err := os.Stat(folder)
	if os.IsNotExist(err) {
		err := os.Mkdir(path.Dir(filePath), 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
