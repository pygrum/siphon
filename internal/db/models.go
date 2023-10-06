package db

import (
	"gorm.io/gorm"
	"time"
)

type Sample struct {
	gorm.Model
	Name      string `json:"name"`
	Path      string `json:"path"`
	FileType  string `json:"file_type"`
	FileSize  uint   `json:"file_size"`
	Signature string `json:"signature"`
	Source    string `json:"source"`

	Hash       string    `json:"hash"`
	UploadTime time.Time `json:"upload_time"`
}

type Agent struct {
	gorm.Model
	AgentID  string `yaml:"agent_id"`
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
	CertPath string `yaml:"certificate_path"`
}
