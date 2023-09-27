package info

import (
	"github.com/pygrum/siphon/internal/commands/samples"
	"github.com/pygrum/siphon/internal/db"
	"github.com/pygrum/siphon/internal/db/models"
	"github.com/pygrum/siphon/internal/logger"
	"strconv"
)

func InfoCmd(ids ...string) {
	ids = clean(ids)
	var sampleList []models.Sample
	for _, id := range ids {
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
		sampleList = append(sampleList, *spl)
	}
	samples.RenderTable(sampleList)
}

func clean(array []string) []string {
	m := make(map[string]string)
	for _, x := range array {
		m[x] = x
	}
	var cleaned []string
	for x, _ := range m {
		cleaned = append(cleaned, x)
	}
	return cleaned
}
