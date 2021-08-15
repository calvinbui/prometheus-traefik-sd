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
			for _, g := range gFiles {
				if g.FilePath == f {
					logger.Debug(fmt.Sprintf("%s will be deleted once grace period is exceeded", f))
					g.Count++
					continue existingFiles
				}
			}

			logger.Info(fmt.Sprintf("%s will be deleted once grace period has exceeded", f))
			gFiles = append(gFiles, GraceFile{FilePath: f, Count: 1})
		}
	}

	for _, g := range gFiles {
		if g.Count >= gracePeriod {
			logger.Info(fmt.Sprintf("Deleting %s as grace period has exceeded", g.FilePath))
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
