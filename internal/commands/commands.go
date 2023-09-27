package commands

import (
	"github.com/pygrum/siphon/internal/commands/exit"
	"github.com/pygrum/siphon/internal/commands/get"
	"github.com/pygrum/siphon/internal/commands/info"
	"github.com/pygrum/siphon/internal/commands/samples"
	"github.com/pygrum/siphon/internal/commands/sources"
	"github.com/pygrum/siphon/internal/commands/sources/new"
	"github.com/spf13/cobra"
)

func Commands() *cobra.Command {
	cmd := &cobra.Command{}

	sourcesCmd := &cobra.Command{
		Use:   "sources",
		Short: "List currently configured threat intelligence sources",
		Run: func(cmd *cobra.Command, args []string) {
			sources.SourcesCmd()
		},
	}

	var newName, newApiKey, newEndpoint string
	newCmd := &cobra.Command{
		Use:   "new",
		Short: "configure a new integration",
		Run: func(cmd *cobra.Command, args []string) {
			new.NewCmd(newName, newApiKey, newEndpoint)
		},
	}
	newCmd.Flags().StringVarP(&newName, "name", "n", "", "name of new source")
	newCmd.Flags().StringVarP(&newApiKey, "api-key", "k", "", "api key for source")
	newCmd.Flags().StringVarP(&newEndpoint, "endpoint", "e", "", "source API endpoint")
	_ = cobra.MarkFlagRequired(newCmd.Flags(), "name")

	sourcesCmd.AddCommand(newCmd)

	var sampleCount string
	samplesCmd := &cobra.Command{
		Use:   "samples",
		Short: "List the latest samples found by Siphon - default 5",
		Run: func(cmd *cobra.Command, args []string) {
			samples.SamplesCmd(sampleCount)
		},
	}
	samplesCmd.Flags().StringVarP(&sampleCount, "count", "c", "5", "number of samples to retrieve")

	var getID string
	var getOut string
	var getPersist bool
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Download an indexed sample from its original source",
		Run: func(cmd *cobra.Command, args []string) {
			get.GetCmd(getID, getOut, getPersist)
		},
	}
	getCmd.Flags().StringVarP(&getID, "id", "", "", "the id or name of the sample to download")
	getCmd.Flags().StringVarP(&getID, "outfile", "o", "", "the save name of the sample")
	getCmd.Flags().BoolVarP(&getPersist, "persist", "p", false, "save the sample to a permanent location")
	_ = cobra.MarkFlagRequired(getCmd.Flags(), "id")

	infoCmd := &cobra.Command{
		Use:   "info [id...]",
		Short: "Get information about 1 or more samples (querying by name or id)",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			info.InfoCmd(args...)
		},
	}

	exitCmd := &cobra.Command{
		Use:   "exit",
		Short: "Exit the application",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			exit.ExitCmd()
		},
	}

	cmd.AddCommand(infoCmd)
	cmd.AddCommand(getCmd)
	cmd.AddCommand(sourcesCmd)
	cmd.AddCommand(samplesCmd)
	cmd.AddCommand(exitCmd)

	cmd.CompletionOptions.HiddenDefaultCmd = true

	return cmd
}
