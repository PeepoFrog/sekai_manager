package mnemonicderiver

// mnemonicsgenerator "github.com/KiraCore/tools/validator-key-gen/MnemonicsGenerator"
import (
	"fmt"

	vlg "github.com/KiraCore/tools/validator-key-gen/MnemonicsGenerator"
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
	// TODO path set as variables or constants
	sekaidConfigFolder := homeFolder + "/config"
	fmt.Println(sekaidConfigFolder)
	var err error
	//creating sekaid home
	// err := os.Mkdir(sekaidHome, 0755)
	// if err != nil {
	// 	if !os.IsExist(err) {
	// 		return fmt.Errorf("unable to create <%s> folder, err: %w", sekaidHome, err)
	// 	}
	// }
	// //creating sekaid's config folder
	// err = os.Mkdir(sekaidConfigFolder, 0755)
	// if err != nil {
	// 	if !os.IsExist(err) {
	// 		return fmt.Errorf("unable to create <%s> folder, err: %w", sekaidConfigFolder, err)
	// 	}
	// }

	err = vlg.GeneratePrivValidatorKeyJson(mnemonicSet.ValidatorValMnemonic, sekaidConfigFolder+"/priv_validator_key.json", vlg.DefaultPrefix, vlg.DefaultPath)
	if err != nil {
		return fmt.Errorf("unable to generate priv_validator_key.json: %w", err)
	}
	err = vlg.GenerateValidatorNodeKeyJson(mnemonicSet.ValidatorNodeMnemonic, sekaidConfigFolder+"/node_key.json", vlg.DefaultPrefix, vlg.DefaultPath)
	if err != nil {
		return fmt.Errorf("unable to generate node_key.json: %w", err)
	}
	return nil
}
