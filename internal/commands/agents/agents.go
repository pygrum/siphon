package agents

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pygrum/siphon/internal/db"
	"github.com/pygrum/siphon/internal/logger"
)

var conn *db.Conn

func init() {
	conn = db.Initialize()
}

func AgentsCmd() {
	// Do nothing if zero sample count
	agents := conn.Agents()
	if len(agents) == 0 {
		logger.Info("no samples loaded - check your internet connection or API configuration")
	} else {
		RenderTable(agents)
	}
}

func RenderTable(agents []db.Agent) {
	t := table.NewWriter()
	tmp := table.Table{}
	tmp.Render()

	t.SetStyle(table.StyleBold)

	header := table.Row{"ID", "NAME", "ENDPOINT", "CERTIFICATE", "CREATION TIME"}
	for _, a := range agents {
		row := table.Row{
			a.AgentID,
			a.Name,
			a.Endpoint,
			a.CertPath,
			a.CreatedAt,
		}
		t.AppendRow(row)
	}
	t.AppendHeader(header)
	t.SetAutoIndex(false)
	fmt.Println(t.Render())
}
