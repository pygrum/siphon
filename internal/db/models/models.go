package models

import (
	"gorm.io/gorm"
	"time"
)

type Sample struct {
	gorm.Model
	Name      string
	FileType  string
	FileSize  uint
	Signature string
	Source    string

	Hash       string
	UploadTime time.Time
}
