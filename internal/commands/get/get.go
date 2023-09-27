package get

import (
	"github.com/pygrum/siphon/internal/db"
	"github.com/pygrum/siphon/internal/db/models"
	"github.com/pygrum/siphon/internal/integrations/malwarebazaar"
	"github.com/pygrum/siphon/internal/logger"
	"github.com/pygrum/siphon/internal/siphon"
	"io"
	"strconv"
)

func GetCmd(id, out string, persist bool) {
	uid, err := strconv.Atoi(id)
	spl := &models.Sample{}
	if err != nil {
		spl = db.SampleByName(id)
	} else {
		spl = db.SampleByID(uint(uid))
	}
	if spl == nil {
		logger.Errorf("sample '%s': not found", id)
		return
	}
	if spl.Source == "MalwareBazaar" {
		fileBytes := getFromMB(spl)
		if fileBytes != nil {
			if len(out) == 0 {
				out = spl.Hash + "." + spl.FileType + ".zip" // MalwareBazaar returns as zip
			} else {
				out += ".zip"
			}
			f, err := siphon.AddFile(out, fileBytes, 0700, persist)
			if err != nil {
				logger.Errorf("failed to save file to %s: %v", out, err)
				return
			}
			logger.Notifyf("saved as archive to registry (%s)", f)
			logger.Notifyf("password: 'infected'")
		}
	}
}

func getFromMB(spl *models.Sample) []byte {
	f := malwarebazaar.NewFetcher()
	if f == nil {
		logger.Error("cannot fetch from MalwareBazaar: not configured correctly")
		return nil
	}
	readCloser, err := f.Download(spl.Hash)
	if err != nil {
		logger.Errorf("download failed: %v", err)
		return nil
	}
	bytes, err := io.ReadAll(readCloser)
	if err != nil {
		logger.Errorf("cannot read response body: %v", err)
	}
	return bytes
}
