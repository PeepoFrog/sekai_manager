package cmd

import (
	"fmt"

	vlg "github.com/KiraCore/tools/validator-key-gen/MnemonicsGenerator"
	mnemonicderiver "github.com/PeepoFrog/sekai_manager/src/instances_manager/mnemonic_deriver"
	"github.com/PeepoFrog/sekai_manager/src/types"
	"github.com/spf13/cobra"
)

// newDeriveValidatorFromMasterCmd is a leaf under root.
func newDeriveValidatorFromMasterCmd(app *types.ManagerConfig) *cobra.Command {
	var (
		mnemonic  string
		path      string
		prefix    string
		outFolder string
	)

	cmd := &cobra.Command{
		Use:     "deriveValidatorFromMaster",
		Aliases: []string{"derive-validator-from-master"},
		Short:   "Derive a validator from a master key",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if mnemonic == "" {
				return fmt.Errorf("mnemonic cannot be empty (use --mnemonic or -m)")
			}
			if outFolder == "" {
				return fmt.Errorf("out folder is required (use --out or -o)")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return deriveMnemonicFromMaster(mnemonic, prefix, path, outFolder)
		},
	}

	// ---- flags ----
	cmd.Flags().StringVarP(&mnemonic, "mnemonic", "m", "", "BIP39 mnemonic (REQUIRED)")
	cmd.Flags().StringVarP(&path, "path", "p", vlg.DefaultPath, "Derivation path (BIP44-style)")
	cmd.Flags().StringVarP(&prefix, "prefix", "x", vlg.DefaultPrefix, "Derivation prefix (BIP44-style)")
	cmd.Flags().StringVarP(&outFolder, "out", "o", "", "Output directory (REQUIRED)")

	// Optional UX sugar
	_ = cmd.MarkFlagRequired("mnemonic")
	_ = cmd.MarkFlagRequired("out")

	return cmd
}

func deriveMnemonicFromMaster(masterMnemonic, prefix, path, outFolder string) error {
	return mnemonicderiver.DeliverMnemonicKeysFromMaster(masterMnemonic, prefix, path, outFolder)
}
