package get

import (
	"github.com/pygrum/siphon/cmd/generator/generator"
	"github.com/pygrum/siphon/internal/db"
	"github.com/pygrum/siphon/internal/integrations/agent"
	"github.com/pygrum/siphon/internal/integrations/malwarebazaar"
	"github.com/pygrum/siphon/internal/logger"
	"github.com/pygrum/siphon/internal/siphon"
	"io"
	"strconv"
	"strings"
)

func GetCmd(id, out string, persist bool) {
	conn := db.Initialize()
	uid, err := strconv.Atoi(id)
	spl := &db.Sample{}
	if err != nil {
		spl = conn.SampleByName(id)
	} else {
		spl = conn.SampleByID(uint(uid))
	}
	if spl == nil {
		logger.Errorf("sample '%s': not found", id)
		return
	}
	var fileBytes []byte
	if strings.ToLower(spl.Source) == malwarebazaar.Source || generator.IsAgentID(spl.Source) {
		if strings.ToLower(spl.Source) == malwarebazaar.Source {
			fileBytes = getFromMB(spl)
		} else if generator.IsAgentID(spl.Source) {
			agt := conn.AgentByID(spl.Source)
			if agt == nil {
				logger.Errorf("agent AgentID '%s' is not present in database", spl.Source)
				return
			}
			fileBytes = getFromAgent(agt, spl)
		}
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

func getFromMB(spl *db.Sample) []byte {
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

func getFromAgent(agt *db.Agent, spl *db.Sample) []byte {
	f := agent.NewFetcher()
	readCloser, err := f.Download(agt, spl.Hash)
	if err != nil {
		logger.Errorf("download failed: %v", err)
		return nil
	}
	defer readCloser.Close()
	bytes, err := io.ReadAll(readCloser)
	if err != nil {
		logger.Errorf("cannot read response body: %v", err)
	}
	return bytes
}
