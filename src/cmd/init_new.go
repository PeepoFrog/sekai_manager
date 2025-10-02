package cmd

import (
	"fmt"

	"github.com/PeepoFrog/sekai_manager/src/types"
	"github.com/spf13/cobra"
)

// newNewCmd is a leaf under init.
func newNewCmd(app *types.ManagerConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "new",
		Short: "Create a new setup",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: implement new logic
			fmt.Println("init new: not implemented yet")
			return nil
		},
	}
}
