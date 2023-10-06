package sources

import (
	"github.com/pygrum/siphon/internal/logger"
	"github.com/spf13/viper"
	"slices"
	"strings"
)

var (
	SupportedSources = []string{
		//"virustotal",
		"malwarebazaar",
		//"hybridanalysis",
		//"virusshare",
	}
)

type Source struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
	ApiKey   string `yaml:"apikey"`
}

func Sources() []Source {
	var sources []Source
	_ = viper.UnmarshalKey("sources", &sources)
	return sources
}

func FindSource(sourceName string) *Source {
	var sources []Source
	_ = viper.UnmarshalKey("sources", &sources)
	if len(sources) == 0 {
		return nil
	}
	for _, s := range sources {
		if strings.ToLower(s.Name) == strings.ToLower(sourceName) {
			return &s
		}
	}
	return nil
}

func SourcesCmd() {
	var sources []Source
	if err := viper.UnmarshalKey("sources", &sources); err != nil {
		logger.Errorf("failed to parse configuration file %s: %v\n", viper.ConfigFileUsed(), err)
		return
	}
	if len(sources) == 0 {
		logger.Info("no sources configured.")
		return
	}
	for _, s := range sources {
		if !slices.Contains(SupportedSources, strings.ToLower(s.Name)) {
			logger.Errorf("'%s' is not a supported integration", s.Name)
			continue
		}
		if len(s.ApiKey) == 0 {
			logger.Warnf("the API key for %s has not been set", s.Name)
			continue
		}
		if len(s.Endpoint) == 0 {
			logger.Warnf("the API endpoint for %s has not been set", s.Name)
			continue
		}
		logger.Notifyf("%s - fully configured", s.Name)
	}
}
