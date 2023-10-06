package samples

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pygrum/siphon/internal/db"
	"github.com/pygrum/siphon/internal/logger"
	"strconv"
	"strings"
)

const (
	truncLimit = 40
)

var conn *db.Conn

func init() {
	conn = db.Initialize()
}

func SamplesCmd(sampleCount string, noTrunc bool) {
	sc := 0
	i, err := strconv.Atoi(sampleCount)
	if err != nil {
		if strings.ToLower(sampleCount) == "all" {
			sc = conn.Count()
		} else {
			logger.Errorf("valid integer argument required")
			return
		}
	} else {
		if i < 0 {
			logger.Errorf("valid integer argument required")
			return
		}
		sc = i
	}
	// Do nothing if zero sample count
	if sc > 0 {
		samples := conn.Samples(sc)
		if len(samples) == 0 {
			logger.Info("no samples loaded - check your internet connection or API configuration")
		} else {
			RenderTable(samples, noTrunc)
		}
	}
}

func RenderTable(samples []db.Sample, v bool) {
	t := table.NewWriter()
	tmp := table.Table{}
	tmp.Render()

	t.SetStyle(table.StyleBold)

	header := table.Row{"ID", "NAME", "TYPE", "SIGNATURE", "HASH", "SIZE", "SOURCE", "UPLOADED AT"}
	for _, s := range samples {
		name := s.Name
		if !v {
			if len(s.Name) > truncLimit {
				name = s.Name[:truncLimit]
				name += "..."
			}
		}
		row := table.Row{
			s.ID,
			name,
			s.FileType,
			s.Signature,
			s.Hash,
			s.FileSize,
			s.Source,
			s.UploadTime.String(),
		}
		t.AppendRow(row)
	}
	t.AppendHeader(header)
	t.SetAutoIndex(false)
	fmt.Println(t.Render())
}
