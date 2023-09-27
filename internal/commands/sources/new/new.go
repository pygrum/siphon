package new

import (
	"github.com/pygrum/siphon/internal/commands/sources"
	"github.com/pygrum/siphon/internal/logger"
	"github.com/spf13/viper"
	"slices"
	"strings"
)

func NewCmd(name, apikey, endpoint string) {
	name = strings.ToLower(name)
	srcs := sources.Sources()
	if !slices.Contains(sources.SupportedSources, name) {
		logger.Errorf("'%s' integration is not supported", name)
		return
	}
	var update bool
	for i, src := range srcs {
		if strings.ToLower(src.Name) == name {
			update = true
			if len(endpoint) != 0 {
				src.Endpoint = endpoint
			}
			if len(apikey) != 0 {
				src.ApiKey = apikey
			}
			srcs[i] = src
			break
		}
	}
	if !update {
		srcs = append(srcs, sources.Source{
			Name:     name,
			ApiKey:   apikey,
			Endpoint: endpoint,
		})
	}
	viper.Set("sources", srcs)
	if err := viper.WriteConfig(); err != nil {
		logger.Errorf("failed to save new configuration: %v", err)
	}
	logger.Notifyf("configuration for %s successfully updated", name)
}
