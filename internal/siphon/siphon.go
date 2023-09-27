package siphon

import (
	"github.com/pygrum/siphon/internal/logger"
	"os"
	"path/filepath"
)

var siphon struct {
	rootDir   string
	sampleDir string
	tmpDir    string
}

func root() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".siphon")
}

func init() {
	siphon.rootDir = root()
	siphon.sampleDir = filepath.Join(siphon.rootDir, "samples")
	siphon.tmpDir = filepath.Join(os.TempDir(), "samples")
}

func fullPath(samplePath string) string {
	return filepath.Join(siphon.sampleDir, samplePath)
}

func fullTmpPath(samplePath string) string {
	return filepath.Join(siphon.tmpDir, samplePath)
}

func AddFile(path string, data []byte, perm os.FileMode, persist bool) (string, error) {
	full := fullTmpPath(path)
	if persist {
		full = fullPath(path)
	}
	if _, err := os.Stat(filepath.Dir(full)); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(full), 0700)
		if err != nil {
			logger.Silentf("failed to create parent folders: %v", err)
			return "", err
		}
	}
	return full, os.WriteFile(full, data, perm)
}
