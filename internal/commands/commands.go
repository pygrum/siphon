package commands

import (
	"github.com/pygrum/siphon/internal/commands/agents"
	"github.com/pygrum/siphon/internal/commands/agents/set"
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
	var sampleNoTruncate bool
	samplesCmd := &cobra.Command{
		Use:   "samples",
		Short: "List the latest samples found by Siphon - default 5",
		Run: func(cmd *cobra.Command, args []string) {
			samples.SamplesCmd(sampleCount, sampleNoTruncate)
		},
	}
	samplesCmd.Flags().StringVarP(&sampleCount, "count", "c", "5", "number of samples to retrieve (use 'all' to retrieve all samples)")
	samplesCmd.Flags().BoolVarP(&sampleNoTruncate, "no-truncate", "v", false, "don't truncate sample names")

	var getOut string
	var getPersist bool
	getCmd := &cobra.Command{
		Use:   "get [id|name]",
		Short: "Download an indexed sample from its original source",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			get.GetCmd(args[0], getOut, getPersist)
		},
	}
	getCmd.Flags().StringVarP(&getOut, "outfile", "o", "", "the save name of the sample")
	getCmd.Flags().BoolVarP(&getPersist, "persist", "p", false, "save the sample to a permanent location")
	_ = cobra.MarkFlagRequired(getCmd.Flags(), "id")

	var infoNoTruncate bool
	infoCmd := &cobra.Command{
		Use:   "info [id...]",
		Short: "Get information about 1 or more samples (querying by name or id)",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			info.InfoCmd(infoNoTruncate, args...)
		},
	}
	infoCmd.Flags().BoolVarP(&infoNoTruncate, "no-truncate", "v", false, "don't truncate sample names")

	exitCmd := &cobra.Command{
		Use:   "exit",
		Short: "Exit the application",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			exit.ExitCmd()
		},
	}

	agentsCmd := &cobra.Command{
		Use:   "agents",
		Short: "View known agents",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			agents.AgentsCmd()
		},
	}

	var endpoint, certFile string
	setCmd := &cobra.Command{
		Use:   "set [id]",
		Short: "Configure agent integration parameters",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			set.SetCmd(args[0], endpoint, certFile)
		},
	}
	setCmd.Flags().StringVarP(&endpoint, "endpoint", "e", "", "API endpoint for the agent")
	setCmd.Flags().StringVarP(&certFile, "cert-file", "c", "", "path to agent certificate file")
	agentsCmd.AddCommand(setCmd)

	cmd.AddCommand(infoCmd)
	cmd.AddCommand(getCmd)
	cmd.AddCommand(sourcesCmd)
	cmd.AddCommand(agentsCmd)
	cmd.AddCommand(samplesCmd)
	cmd.AddCommand(exitCmd)

	cmd.CompletionOptions.HiddenDefaultCmd = true

	return cmd
}
