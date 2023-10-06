package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/pygrum/siphon/internal/db"
)

type Agent struct {
	Router         *gin.Engine
	ID             string
	BindAddress    string
	CertFile       string
	KeyFile        string
	ClientCertData []byte // Get this at compile time
	Conn           *db.AgentConn
}

type SampleResponse struct {
	Status string      `json:"status"`
	Data   []db.Sample `json:"data"`
}
