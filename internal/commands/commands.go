package commands

import (
	"github.com/pygrum/siphon/internal/commands/exit"
	"github.com/pygrum/siphon/internal/commands/get"
	"github.com/pygrum/siphon/internal/commands/samples"
	"github.com/pygrum/siphon/internal/commands/sources"
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

	exitCmd := &cobra.Command{
		Use:   "exit",
		Short: "Exit the application",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			exit.ExitCmd()
		},
	}

	cmd.AddCommand(getCmd)
	cmd.AddCommand(sourcesCmd)
	cmd.AddCommand(samplesCmd)
	cmd.AddCommand(exitCmd)

	cmd.CompletionOptions.HiddenDefaultCmd = true

	return cmd
}
