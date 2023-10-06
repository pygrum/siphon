package main

import (
	"fmt"
	"github.com/pygrum/siphon/cmd/generator/generator"
	"github.com/pygrum/siphon/internal/logger"
	"github.com/pygrum/siphon/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	fmt.Printf("{_/¬ SIPHON GENERATOR %s ¬\\_}\n\n}", version.VersionString())
}

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "generator",
		Short: "A utility for Siphon agent generation",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			viper.SetConfigFile(cfgFile)
			if err := viper.ReadInConfig(); err != nil {
				logger.Fatalf("reading configuration file failed: %v", err)
			}
			if err := generator.Generate(); err != nil {
				logger.Fatal(err)
			}
			logger.Notifyf("agent has successfully been built. For installation instructions, see the docs: %s",
				"https://github.com/pygrum/siphon/blob/main/docs/DOCS.md",
			)
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "generator configuration file - see https://github.com/pygrum/siphon/blob/main/docs/DOCS.md for help")
	_ = cobra.MarkFlagRequired(rootCmd.PersistentFlags(), "config")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}
