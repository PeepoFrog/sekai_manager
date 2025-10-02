package cmd

import (
	"github.com/PeepoFrog/sekai_manager/src/types"
	"github.com/spf13/cobra"
)

// NewRootCmd constructs the root command and wires subcommands.
func NewRootCmd(app *types.ManagerConfig) *cobra.Command {
	root := &cobra.Command{
		Use:   "app",
		Short: "CLI root command",
		Long:  "An example CLI showing a subcommand tree with init/{join,new}, deriveValidatorFromMaster, and status.",
	}

	// Attach subcommands
	root.AddCommand(newInitCmd(app))
	root.AddCommand(newDeriveValidatorFromMasterCmd(app))
	root.AddCommand(newStatusCmd(app))

	return root
}
