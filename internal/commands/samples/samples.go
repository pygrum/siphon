package samples

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pygrum/siphon/internal/db"
	"github.com/pygrum/siphon/internal/db/models"
	"github.com/pygrum/siphon/internal/logger"
	"strconv"
	"strings"
)

func SamplesCmd(sampleCount string) {
	sc := 0
	i, err := strconv.Atoi(sampleCount)
	if err != nil {
		if strings.ToLower(sampleCount) == "all" {
			sc = db.Count()
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
		samples := db.Samples(sc)
		if len(samples) == 0 {
			logger.Info("no samples loaded - check your internet connection or API configuration")
		} else {
			RenderTable(samples)
		}
	}
}

func RenderTable(samples []models.Sample) {
	t := table.NewWriter()
	tmp := table.Table{}
	tmp.Render()

	t.SetStyle(table.StyleBold)

	header := table.Row{"ID", "NAME", "TYPE", "SIGNATURE", "HASH", "SIZE", "SOURCE", "UPLOADED AT"}
	for _, s := range samples {
		row := table.Row{
			s.ID,
			s.Name,
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
