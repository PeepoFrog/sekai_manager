package cmd

import (
	"fmt"

	"github.com/PeepoFrog/sekai_manager/src/types"
	"github.com/spf13/cobra"
)

// newJoinCmd is a leaf under init.
func newJoinCmd(app *types.ManagerConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "join",
		Short: "Join an existing setup",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: implement join logic
			fmt.Println("init join: not implemented yet")
			return nil
		},
	}
}
