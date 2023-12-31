package integrations

import (
	"github.com/pygrum/siphon/internal/integrations/agent"
	"github.com/pygrum/siphon/internal/integrations/malwarebazaar"
	"github.com/pygrum/siphon/internal/logger"
	"github.com/spf13/viper"
	"time"
)

func Refresh() {
	r := viper.GetInt("RefreshRate")
	if r < 1 {
		logger.Silentf("invalid configuration: refresh rate must be 1 minute or more")
	}
	ticker := time.NewTicker(time.Duration(r) * time.Minute)
	for range ticker.C {
		mbFetcher := malwarebazaar.NewFetcher()
		if mbFetcher != nil {
			go mbFetcher.GetRecent()
		}
		agFetcher := agent.NewFetcher()
		go agFetcher.GetRecent()
	}
}
