package agent

import (
	"github.com/pygrum/siphon/internal/agent/controllers"
	"github.com/pygrum/siphon/internal/agent/monitor"
	"github.com/pygrum/siphon/internal/logger"
)

func Initialize(id, iFace, port, clientCertData string) {
	agent := controllers.NewAgent(id, iFace, port, clientCertData)
	// person who compiled the agent is the only one who can add and remove users - their username 'root' is special
	setupControllers(agent)
	scout, err := monitor.NewScout(agent)
	if err != nil {
		logger.Fatalf("failed to create a new scout: %v", err)
	}
	if err := scout.Start(); err != nil {
		logger.Fatalf("failed to start scout: %v", err)
	}
	logger.Fatal(scout.RunTLS())
}

func setupControllers(agent *controllers.Agent) {
	apiRouter := agent.Router.Group("/api")
	{
		apiRouter.GET("samples", agent.GetSamples)
		apiRouter.GET("download", agent.GetSampleByHash)
	}
}
