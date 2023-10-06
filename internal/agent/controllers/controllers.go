package controllers

import (
	"crypto/tls"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/pygrum/siphon/internal/agent/utils"
	"github.com/pygrum/siphon/internal/db"
	"github.com/pygrum/siphon/internal/logger"
	"net"
	"net/http"
	"os"
	"time"
)

// TODO: Use gin to setup API endpoints, that will be used to (DONE) info about recent samples, and (DONE) samples by hash.

const (
	StatusOk       = "ok"
	StatusNotFound = "not_found"
)

func NewAgent(id, iFace, port, clientCertData string) *Agent {
	certData, err := base64.StdEncoding.DecodeString(clientCertData)
	if err != nil {
		logger.Fatalf("could not decode client certificate: %v", err)
	}
	a := &Agent{
		Router:         gin.Default(),
		ID:             id,
		ClientCertData: certData,
		BindAddress:    net.JoinHostPort(iFace, port),
		Conn:           db.AgentInitialize(),
	}
	return a
}

func (a *Agent) RunTLS(tlsConfig *tls.Config) error {
	server := &http.Server{
		Addr:      a.BindAddress,
		Handler:   a.Router,
		TLSConfig: tlsConfig,
	}
	return server.ListenAndServeTLS(a.CertFile, a.KeyFile)
}

func (a *Agent) GetSamples(c *gin.Context) {
	// Return samples discovered within the last hour
	samples, err := a.Conn.SamplesByTime(time.Now().Add(-1 * time.Hour))
	if err != nil {
		c.JSON(http.StatusInternalServerError, SampleResponse{
			Status: err.Error(),
		})
	}
	c.JSON(http.StatusOK, SampleResponse{
		Status: StatusOk,
		Data:   samples,
	})
}

func (a *Agent) GetSampleByHash(c *gin.Context) {
	hash := c.Query("sha256_hash")
	// Return samples discovered within the last hour
	sample, err := a.Conn.SampleByHash(hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, SampleResponse{
			Status: err.Error(),
		})
		return
	}
	if sample == nil {
		c.JSON(http.StatusOK, SampleResponse{
			Status: StatusNotFound,
		})
		return
	}
	file, err := utils.ZipFile(sample, "infected")
	if err != nil {
		c.JSON(http.StatusOK, SampleResponse{
			Status: StatusNotFound,
		})
		return
	}
	bytes, err := os.ReadFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, SampleResponse{
			Status: err.Error(),
		})
		return
	}
	c.Data(http.StatusOK, "application/zip", bytes)
	_ = os.Remove(file)
}
