package cmd

import (
	"fmt"

	"github.com/PeepoFrog/sekai_manager/src/types"
	"github.com/spf13/cobra"
)

// newStatusCmd is a leaf under root.
func newStatusCmd(app *types.ManagerConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current status (stub)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: implement status logic
			fmt.Println("status: OK (stub)")
			return nil
		},
	}
}
