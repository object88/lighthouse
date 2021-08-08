package cmd

import (
	"time"

	"github.com/object88/lighthouse/apps/collector/cmd/run"
	"github.com/object88/lighthouse/internal/cmd/common"
	"github.com/object88/lighthouse/internal/cmd/completion"
	"github.com/object88/lighthouse/internal/cmd/version"
	"github.com/spf13/cobra"
)

// const bashCompletionFunc = `
// __lighthouse_get_outputs()
// {
// 	COMPREPLY=( "json", "json-compressed", "text" )
// }
// `

// InitializeCommands sets up the cobra commands
func InitializeCommands() *cobra.Command {
	ca, rootCmd := createRootCommand()

	rootCmd.AddCommand(
		completion.CreateCommand(ca),
		run.CreateCommand(ca),
		version.CreateCommand(ca),
	)

	return rootCmd
}

func createRootCommand() (*common.CommonArgs, *cobra.Command) {
	ca := common.NewCommonArgs()

	var start time.Time
	cmd := &cobra.Command{
		Use:   "lighthouse-collector",
		Short: "lighthouse-collector monitors a helm installation, upgrade, or deletion",
		// BashCompletionFunction: bashCompletionFunc,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			start = time.Now()
			ca.Evaluate()

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
		PersistentPostRunE: func(cmd *cobra.Command, _ []string) error {
			ca.ReportDuration(cmd, start)
			return nil
		},
	}

	flags := cmd.PersistentFlags()
	ca.Setup(flags)

	return ca, common.TraverseRunHooks(cmd)
}
