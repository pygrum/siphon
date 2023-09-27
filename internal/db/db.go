package db

import (
	"github.com/pygrum/siphon/internal/db/models"
	"github.com/pygrum/siphon/internal/logger"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	l "gorm.io/gorm/logger"
	"sync"

	"os"
	"path/filepath"
)

var conn struct {
	File string
	DB   *gorm.DB
}

var m sync.Mutex

func Initialize() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	dbFile := filepath.Join(home, ".siphon", "siphon.db")
	if _, err = os.Stat(dbFile); os.IsNotExist(err) {
		err := os.WriteFile(dbFile, nil, 0666)
		if err != nil {
			logger.Fatalf("could not initialize database: %v", err)
		}
	}

	conn.File = dbFile
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{
		Logger: l.Default.LogMode(l.Silent),
	})
	if err != nil {
		logger.Fatalf("failed to open database at %s: %v", dbFile, err)
	}
	_ = db.AutoMigrate(models.Sample{})
	conn.DB = db
}

func Samples(count int) []models.Sample {
	var samples []models.Sample
	conn.DB.Order("upload_time DESC, created_at DESC").Limit(count).Find(&samples)
	return samples
}

func Count() int {
	var samples []models.Sample
	var count int64
	conn.DB.Find(&samples).Count(&count)
	return int(count)
}

func SampleByHash(md5Hash string) *models.Sample {
	sample := &models.Sample{}
	m.Lock()
	defer m.Unlock()
	conn.DB.Where("hash = ?", md5Hash).First(&sample)
	// Return if empty sample received
	if (models.Sample{}) == *sample {
		return nil
	}
	return sample
}

func SampleByName(name string) *models.Sample {
	sample := &models.Sample{}
	m.Lock()
	defer m.Unlock()
	conn.DB.Where("name = ?", name).First(&sample)
	// Return if empty sample received
	if (models.Sample{}) == *sample {
		return nil
	}
	return sample
}

func SampleByID(id uint) *models.Sample {
	sample := &models.Sample{}
	m.Lock()
	defer m.Unlock()
	conn.DB.Where("id = ?", id).First(&sample)
	// Return if empty sample received
	if (models.Sample{}) == *sample {
		return nil
	}
	return sample
}

func AddSample(sample *models.Sample) error {
	m.Lock()
	defer m.Unlock()
	result := conn.DB.Create(sample)
	return result.Error
}
