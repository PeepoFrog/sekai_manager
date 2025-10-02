package cmd

import (
	"fmt"

	"github.com/PeepoFrog/sekai_manager/src/types"
	"github.com/spf13/cobra"
)

// newDeriveValidatorFromMasterCmd is a leaf under root.
func newDeriveValidatorFromMasterCmd(app *types.ManagerConfig) *cobra.Command {
	return &cobra.Command{
		Use:     "deriveValidatorFromMaster",
		Aliases: []string{"derive-validator-from-master"},
		Short:   "Derive a validator from a master key (stub)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: implement derivation logic
			fmt.Println("deriveValidatorFromMaster: not implemented yet")
			return nil
		},
	}
}
