package helpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/calvinbui/prometheus-traefik-sd/internal/logger"
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
			for i, g := range gFiles {
				if g.FilePath == f {
					gFiles[i].Count++
					logger.Debug(fmt.Sprintf("%s grace period count is now %v", g.FilePath, g.Count))
					continue existingFiles
				}
			}

			logger.Info(fmt.Sprintf("%s will be deleted once grace period has exceeded", f))
			gFiles = append(gFiles, GraceFile{FilePath: f, Count: 1})
		}
	}

	logger.Debug(fmt.Sprintf("Grace period count is set to %v", gracePeriod))
	var newGFiles []GraceFile
	for _, g := range gFiles {
		if g.Count >= gracePeriod {
			logger.Info(fmt.Sprintf("Deleting %s as grace period has exceeded (%v)", g.FilePath, g.Count))
			if err = os.Remove(g.FilePath); err != nil {
				return nil, err
			}
		} else {
			newGFiles = append(newGFiles, g)
		}
	}

	return newGFiles, nil
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
