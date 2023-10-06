package db

import (
	"encoding/json"
	"github.com/pygrum/siphon/internal/logger"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	l "gorm.io/gorm/logger"
	"runtime"
	"strings"
	"sync"
	"time"

	"os"
	"path/filepath"
)

type Conn struct {
	File string
	DB   *gorm.DB
}

type AgentConn struct {
	File string
}

var m sync.Mutex

func Initialize() *Conn {
	conn := Conn{}
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
	_ = db.AutoMigrate(Sample{}, Agent{})
	conn.DB = db
	return &conn
}

func AgentInitialize() *AgentConn {
	conn := AgentConn{}
	// Safe path to store siphon data based on OS
	var RestrictedPath string
	switch runtime.GOOS {
	case "windows":
		RestrictedPath = "C:\\Windows\\System32\\config"
	default:
		RestrictedPath = "/root"
	}
	if err := os.Mkdir(filepath.Join(RestrictedPath, ".siphon_agent"), 0700); err != nil && !os.IsExist(err) {
		logger.Fatalf("could not create a protected folder in %s: %v", RestrictedPath, err)
	}
	jsonFile := filepath.Join(RestrictedPath, ".siphon_agent", "agent.json")
	if err := os.WriteFile(jsonFile, []byte("[]"), 0600); err != nil {
		logger.Fatalf("could not initialise json database: %v", err)
	}
	conn.File = jsonFile
	return &conn
}

func (conn *Conn) SamplesByTime(dateTime time.Time) []Sample {
	var samples []Sample
	conn.DB.Where("created_at > ?", dateTime.Format(time.DateTime)).Find(&samples)
	return samples
}
func (conn *Conn) Samples(count int) []Sample {
	var samples []Sample
	conn.DB.Order("upload_time DESC, created_at DESC").Limit(count).Find(&samples)
	return samples
}

func (conn *Conn) Agents() []Agent {
	var agents []Agent
	conn.DB.Order("created_at DESC").Find(&agents)
	return agents
}

func (conn *Conn) Count() int {
	var samples []Sample
	var count int64
	conn.DB.Find(&samples).Count(&count)
	return int(count)
}

func (conn *Conn) SampleByHash(sha256hash string) *Sample {
	sample := &Sample{}
	m.Lock()
	defer m.Unlock()
	conn.DB.Where("hash = ?", sha256hash).First(&sample)
	// Return if empty sample received
	if (Sample{}) == *sample {
		return nil
	}
	return sample
}

func (conn *Conn) SampleByName(name string) *Sample {
	sample := &Sample{}
	m.Lock()
	defer m.Unlock()
	conn.DB.Where("name = ?", name).First(&sample)
	// Return if empty sample received
	if (Sample{}) == *sample {
		return nil
	}
	return sample
}

func (conn *Conn) AgentByID(id string) *Agent {
	agent := &Agent{}
	m.Lock()
	defer m.Unlock()
	conn.DB.Where("agent_id = ?", id).First(&agent)
	// Return if empty agent received
	if (Agent{}) == *agent {
		return nil
	}
	return agent
}

func (conn *Conn) AgentByName(name string) *Agent {
	agent := &Agent{}
	m.Lock()
	defer m.Unlock()
	conn.DB.Where("name = ?", name).First(&agent)
	// Return if empty agent received
	if (Agent{}) == *agent {
		return nil
	}
	return agent
}

func (conn *Conn) SampleByID(id uint) *Sample {
	sample := &Sample{}
	m.Lock()
	defer m.Unlock()
	conn.DB.Where("id = ?", id).First(sample)
	// Return if empty sample received
	if (Sample{}) == *sample {
		return nil
	}
	return sample
}

func (conn *Conn) Add(v interface{}) error {
	m.Lock()
	defer m.Unlock()
	result := conn.DB.Save(v)
	return result.Error
}

func (aConn *AgentConn) Add(s *Sample) error {
	samples, err := aConn.Samples()
	if err != nil {
		return err
	}
	samples = append(samples, *s)
	return aConn.Save(samples)
}

func (aConn *AgentConn) Save(samples []Sample) error {
	bytes, err := json.Marshal(samples)
	if err != nil {
		return err
	}
	return os.WriteFile(aConn.File, bytes, 0600)
}

func (aConn *AgentConn) Samples() ([]Sample, error) {
	bytes, err := os.ReadFile(aConn.File)
	if err != nil {
		return nil, err
	}
	var samples []Sample
	if err = json.Unmarshal(bytes, &samples); err != nil {
		return nil, err
	}
	return samples, nil
}

func (aConn *AgentConn) SamplesByTime(dateTime time.Time) ([]Sample, error) {
	samples, err := aConn.Samples()
	if err != nil {
		return nil, err
	}
	var recentSamples []Sample
	for _, sample := range samples {
		// If it was created after an hour ago
		if sample.CreatedAt.After(dateTime) {
			recentSamples = append(recentSamples, sample)
		}
	}
	return recentSamples, nil
}

func (aConn *AgentConn) SampleByHash(sha256hash string) (*Sample, error) {
	samples, err := aConn.Samples()
	if err != nil {
		return nil, err
	}
	for _, sample := range samples {
		if strings.EqualFold(sample.Hash, sha256hash) {
			return &sample, nil
		}
	}
	return nil, nil
}
