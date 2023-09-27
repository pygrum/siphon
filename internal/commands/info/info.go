package info

import (
	"github.com/pygrum/siphon/internal/commands/samples"
	"github.com/pygrum/siphon/internal/db"
	"github.com/pygrum/siphon/internal/db/models"
	"github.com/pygrum/siphon/internal/logger"
	"strconv"
)

func InfoCmd(ids ...string) {
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
	sampleList = clean(sampleList)
	samples.RenderTable(sampleList)
}

func clean(spls []models.Sample) []models.Sample {
	m := make(map[uint]models.Sample)
	for _, x := range spls {
		m[x.ID] = x
	}
	var cleaned []models.Sample
	for _, v := range m {
		cleaned = append(cleaned, v)
	}
	return cleaned
}
