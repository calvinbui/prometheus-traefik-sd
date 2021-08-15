package helpers

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type GraceFile struct {
	FilePath string
	Count    int
}

func DeleteOldTargets(tgs []PromTargetFile, folder string, gFiles []GraceFile, gracePeriod int) ([]GraceFile, error) {
	existingFiles, err := getAllJSONInDirectory(folder)
	if err != nil {
		return nil, err
	}

existingFiles:
	for _, f := range existingFiles {
		if !targetExists(tgs, f) {
			for _, g := range gFiles {
				if g.FilePath == f {
					g.Count++
					continue existingFiles
				}
			}

			gFiles = append(gFiles, GraceFile{FilePath: f, Count: 1})
		}
	}

	for _, g := range gFiles {
		if g.Count >= gracePeriod {
			if err = os.Remove(g.FilePath); err != nil {
				return nil, err
			}
		}
	}

	return gFiles, nil
}

func getAllJSONInDirectory(dir string) ([]string, error) {
	fi, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string

	for _, f := range fi {
		if strings.HasSuffix(f.Name(), ".json") && !f.IsDir() {
			files = append(files, path.Join(dir, f.Name()))
		}
	}

	return files, nil
}

func targetExists(tgs []PromTargetFile, jsonFile string) bool {
	for _, t := range tgs {
		if t.FilePath == jsonFile {
			return true
		}
	}

	return false
}
