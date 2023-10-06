package set

import (
	"github.com/pygrum/siphon/internal/db"
	"github.com/pygrum/siphon/internal/logger"
)

func SetCmd(name, endpoint, certPath string) {
	conn := db.Initialize()
	a := conn.AgentByName(name)
	if a == nil {
		a = conn.AgentByID(name)
		if a == nil {
			logger.Errorf("agent with name/id %s does not exist", name)
			return
		}
	}
	if len(endpoint) != 0 {
		a.Endpoint = endpoint
	}
	if len(certPath) != 0 {
		a.CertPath = certPath
	}
	err := conn.Add(a)
	if err != nil {
		logger.Errorf("failed to update agent fields: %v", err)
		return
	}
	logger.Notify("agent successfully updated")
}
