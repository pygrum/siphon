package console

import (
	"github.com/pygrum/siphon/internal/commands"
	"github.com/pygrum/siphon/internal/integrations"
	"github.com/reeflective/console"
)

type Siphon struct {
	Console *console.Console
}

func Start() {
	var s = &Siphon{
		Console: console.New("siphon"),
	}
	s.setup()
	go s.refresh()
	_ = s.Console.Start()
}

func (s *Siphon) setup() {
	mainMenu := s.Console.ActiveMenu()
	mainMenu.SetCommands(commands.Commands)
	prompt := mainMenu.Prompt()
	prompt.Primary = func() string {
		return "[siphon] "
	}
}

// refresh refreshes sample database every interval
func (s *Siphon) refresh() {
	integrations.Refresh()
	return
}
