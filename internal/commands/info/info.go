package info

import (
	"github.com/pygrum/siphon/internal/commands/samples"
	"github.com/pygrum/siphon/internal/db"
	"github.com/pygrum/siphon/internal/logger"
	"strconv"
)

var conn *db.Conn

func init() {
	conn = db.Initialize()
}

func InfoCmd(noTrunc bool, ids ...string) {
	var sampleList []db.Sample
	for _, id := range ids {
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
		sampleList = append(sampleList, *spl)
	}
	sampleList = clean(sampleList)
	samples.RenderTable(sampleList, noTrunc)
}

func clean(spls []db.Sample) []db.Sample {
	m := make(map[uint]db.Sample)
	for _, x := range spls {
		m[x.ID] = x
	}
	var cleaned []db.Sample
	for _, v := range m {
		cleaned = append(cleaned, v)
	}
	return cleaned
}
