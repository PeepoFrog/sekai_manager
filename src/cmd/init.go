package cmd

import (
	"github.com/PeepoFrog/sekai_manager/src/types"
	"github.com/spf13/cobra"
)

// newInitCmd returns the "init" parent command and adds its leaf subcommands.
func newInitCmd(app *types.ManagerConfig) *cobra.Command {
	c := &cobra.Command{
		Use:   "init",
		Short: "Initialization tasks",
		Long:  "Initialize a setup. Use one of the leaf subcommands: join or new.",
	}

	// Leaf commands
	c.AddCommand(newJoinCmd(app))
	c.AddCommand(newNewCmd(app))
	return c
}
