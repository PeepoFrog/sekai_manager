package cmd

import (
	"fmt"
	"os"

	vlg "github.com/KiraCore/tools/validator-key-gen/MnemonicsGenerator"
	mnemonicderiver "github.com/PeepoFrog/sekai_manager/src/instances_manager/mnemonic_deriver"
	"github.com/PeepoFrog/sekai_manager/src/types"
	"github.com/spf13/cobra"
)

// newDeriveValidatorFromMasterCmd is a leaf under root.
func newDeriveValidatorFromMasterCmd(app *types.ManagerConfig) *cobra.Command {
	// flag storage (lives as long as the command does)
	var (
		mnemonic string
		// words     int
		// folder string
		path      string
		prefix    string
		outFolder string
		// overwrite bool
	)

	cmd := &cobra.Command{
		Use:     "deriveValidatorFromMaster",
		Aliases: []string{"derive-validator-from-master"},
		Short:   "Derive a validator from a master key",
		// Validate before RunE if you want early, clean errors
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if outFolder == "" {
				return fmt.Errorf("out folder is req")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			return deriveMnemonicFromMaster(mnemonic, prefix, path, outFolder)
		},
	}

	// ---- flags ----
	cmd.Flags().StringVarP(&mnemonic, "mnemonic", "m", "", "BIP39 mnemonic to use; if empty, a new one is generated")
	// cmd.Flags().IntVarP(&words, "words", "w", 24, "Word count for a new mnemonic (12,15,18,21,24)")
	cmd.Flags().StringVarP(&path, "path", "pa", vlg.DefaultPath, "Derivation path (BIP44-style)")
	cmd.Flags().StringVarP(&prefix, "prefix", "pr", vlg.DefaultPrefix, "Derivation prefix (BIP44-style)")

	cmd.Flags().StringVarP(&outFolder, "out", "o", "", "Write result to file (default: stdout)")
	// cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Allow overwriting the output file if it exists")

	// examples of extra UX:
	// _ = cmd.MarkFlagFilename("out")           // shell completion hints for files
	// _ = cmd.MarkFlagRequired("path")          // if you want to force a path

	return cmd
}

func deriveMnemonicFromMaster(masterMnemonic, prefix, path, outFolder string) error {
	set, err := mnemonicderiver.GenerateMnemonicsFromMaster(masterMnemonic, prefix, path)
	if err != nil {
		return err
	}
	err = os.MkdirAll(outFolder, 0755)
	if err != nil {
		return err
	}

	err = mnemonicderiver.SetSekaidPrivKeys(set, outFolder)
	if err != nil {
		return err
	}

	return nil
}
