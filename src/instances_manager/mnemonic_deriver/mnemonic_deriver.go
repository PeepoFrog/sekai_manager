package mnemonicderiver

// mnemonicsgenerator "github.com/KiraCore/tools/validator-key-gen/MnemonicsGenerator"
import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	vlg "github.com/KiraCore/tools/validator-key-gen/MnemonicsGenerator"
	"github.com/cosmos/go-bip39"
)

func GenerateMnemonicsFromMaster(masterMnemonic, prefix, path string) (*vlg.MasterMnemonicSet, error) {
	// defaultPrefix := vlg.DefaultPrefix
	// defaultPath := vlg.DefaultPath

	mnemonicSet, err := vlg.MasterKeysGen([]byte(masterMnemonic), prefix, path, "")
	if err != nil {
		return nil, err
	}

	return &mnemonicSet, nil
}

func SetSekaidPrivKeys(mnemonicSet *vlg.MasterMnemonicSet, homeFolder string) error {
	sekaidConfigFolder := filepath.Join(homeFolder, "config")

	// 🔐 Create config dir (secrets → 0700)
	if err := os.MkdirAll(sekaidConfigFolder, 0o700); err != nil {
		return fmt.Errorf("unable to create config dir: %w", err)
	}

	// priv_validator_key.json
	if err := vlg.GeneratePrivValidatorKeyJson(
		mnemonicSet.ValidatorValMnemonic,
		filepath.Join(sekaidConfigFolder, "priv_validator_key.json"),
		vlg.DefaultPrefix,
		vlg.DefaultPath,
	); err != nil {
		return fmt.Errorf("unable to generate priv_validator_key.json: %w", err)
	}

	// node_key.json
	if err := vlg.GenerateValidatorNodeKeyJson(
		mnemonicSet.ValidatorNodeMnemonic,
		filepath.Join(sekaidConfigFolder, "node_key.json"),
		vlg.DefaultPrefix,
		vlg.DefaultPath,
	); err != nil {
		return fmt.Errorf("unable to generate node_key.json: %w", err)
	}

	return nil
}

func DeliverMnemonicKeysFromMaster(masterMnemonic, prefix, path, outFolder string) error {
	valid, invalidWords := CheckMnemonic(masterMnemonic)
	if !valid {
		return fmt.Errorf("invalid mnemonic, invalid words: %v", invalidWords)
	}
	set, err := GenerateMnemonicsFromMaster(masterMnemonic, prefix, path)
	if err != nil {
		return err
	}
	err = os.MkdirAll(outFolder, 0755)
	if err != nil {
		return err
	}

	err = SetSekaidPrivKeys(set, outFolder)
	if err != nil {
		return err
	}

	// 	type MasterMnemonicSet struct {
	//     ValidatorAddrMnemonic []byte
	//     ValidatorValMnemonic  []byte
	//     SignerAddrMnemonic    []byte
	//     ValidatorNodeMnemonic []byte
	//     ValidatorNodeId       []byte
	//     PrivKeyMnemonic       []byte

	stringset := fmt.Sprintf(`
valAddrMnemonic=%s
valValMnemonic=%s
signerAddrMnemonic=%s
valNodeMnemonic=%s
valNodeID=%s`, (set.ValidatorAddrMnemonic), set.ValidatorValMnemonic, set.SignerAddrMnemonic, set.ValidatorNodeMnemonic, set.ValidatorNodeId)
	setFile := filepath.Join(outFolder, "masterSet.txt")
	err = os.WriteFile(setFile, []byte(stringset), 0600)
	return nil
}

// CheckMnemonic prints invalid words (not in BIP39 wordlist) and returns whether the mnemonic is valid.
// Note: If all words are valid but checksum/word-count is wrong, invalidWords will be empty but valid=false.
func CheckMnemonic(mnemonic string) (valid bool, invalidWords []string) {
	words := strings.Fields(strings.ToLower(mnemonic))
	if len(words) == 0 {
		fmt.Println("mnemonic is empty")
		return false, nil
	}

	// Build a set from the currently configured wordlist (default: English).
	// cosmos/go-bip39 exposes WordList as a variable. :contentReference[oaicite:1]{index=1}
	wordSet := make(map[string]struct{}, len(bip39.WordList))
	for _, w := range bip39.WordList {
		wordSet[w] = struct{}{}
	}

	seen := make(map[string]struct{})
	for _, w := range words {
		if _, ok := wordSet[w]; !ok {
			// Keep unique invalid words (easier to read). If you want duplicates/positions, see note below.
			if _, already := seen[w]; !already {
				invalidWords = append(invalidWords, w)
				seen[w] = struct{}{}
			}
		}
	}

	if len(invalidWords) > 0 {
		fmt.Printf("invalid mnemonic words (not in BIP39 wordlist): %v\n", invalidWords)
		return false, invalidWords
	}

	// All words exist in the wordlist; now validate full mnemonic (count + checksum).
	// cosmos/go-bip39 provides IsMnemonicValid. :contentReference[oaicite:2]{index=2}
	if !bip39.IsMnemonicValid(strings.Join(words, " ")) {
		fmt.Println("all words are in the BIP39 wordlist, but the mnemonic is still invalid (wrong word count and/or checksum).")
		return false, nil
	}

	return true, nil
}
