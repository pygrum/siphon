package main

import (
	"errors"
	"fmt"
	"github.com/pygrum/siphon/internal/console"
	"github.com/pygrum/siphon/internal/db"
	"github.com/pygrum/siphon/internal/logger"
	"github.com/pygrum/siphon/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "siphon",
		Short: "Siphon - A CLI-based Threat Intelligence and Asset Feed",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			// Clear screen
			fmt.Println("\033[2J")

			fmt.Println(strings.ReplaceAll(
				title(),
				"{VER}",
				version.VersionString()))
			console.Start()
		},
	}
)

func init() {
	initCfg()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "siphon configuration file")
}

func initCfg() {
	if len(cfgFile) != 0 {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		cfgDir := filepath.Join(home, ".siphon")
		err = os.Mkdir(cfgDir, 0700)
		if !errors.Is(err, os.ErrExist) {
			cobra.CheckErr(err)
		}
		if _, err := os.Stat(filepath.Join(cfgDir, ".siphon.yaml")); os.IsNotExist(err) {
			err = os.WriteFile(filepath.Join(cfgDir, ".siphon.yaml"), cfgBoilerplate(), 0666)
			cobra.CheckErr(err)
		}

		viper.AddConfigPath(cfgDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".siphon")
	}

	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("reading configuration file failed: %v", err)
	}
	db.Initialize()
}

func title() string {
	return `
     ⠀⠀⠀⠀⣀⣤⠶⠒⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣤⣴⠛⢦
     ⠀⣠⣴⠟⠋⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣴⠿⠋⠁⠙⠒⠋
     ⣰⡟⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣴⠟⠁⠀⠀⠀⠀⠀⠀
     ⣿⡁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⡾⠋⠁⠀⠀⠀⠀⠀⠀⠀⠀
     ⠸⣧⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣴⠾⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀     ⠙⢷⣤⡀⠀⠀⠀⠀⠀⠀⠀⣀⣤⡾⠟⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀     ⠀⠀⠉⠻⠷⢶⣶⣶⣶⠶⠟⠛⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
                  
       SIPHON - MALWARE FEED {VER}
     https://github.com/pygrum/siphon
`
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}

func cfgBoilerplate() []byte {
	return []byte(`refreshrate: 5 # Refresh sample list every 5 minutes
Sources:
- name: MalwareBazaar
  endpoint: null
  apikey: null
`)
}
