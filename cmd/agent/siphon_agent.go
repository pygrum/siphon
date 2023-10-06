package main

import (
	"github.com/pygrum/siphon/internal/agent"
	"github.com/pygrum/siphon/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	AgentID        string
	Interface      string
	Port           string
	ClientCertData string // Base64 Encoded certificate data
	cfgFile        string

	rootCmd = &cobra.Command{
		Use:   "siphon_agent",
		Short: "A Honeypot-Resident Sample Curator",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			viper.SetConfigFile(cfgFile)
			if err := viper.ReadInConfig(); err != nil {
				logger.Fatalf("reading configuration file failed: %v", err)
			}
			agent.Initialize(AgentID, Interface, Port, ClientCertData)
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "agent configuration file")
	_ = cobra.MarkFlagRequired(rootCmd.PersistentFlags(), "config")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}
