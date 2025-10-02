package cfg

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/PeepoFrog/sekai_manager/src/types"
	"github.com/pelletier/go-toml/v2"
)

const (
	MANAGER_HOME_FOLDER_NAME string = ".sekaid_manager"
	MANAGER_CONFIG_FILE_NAME string = "cfg.toml"
)

func DefaultCfg() (*types.ManagerConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	homePath := filepath.Join(homeDir, MANAGER_HOME_FOLDER_NAME)
	cfgPath := filepath.Join(homePath, MANAGER_CONFIG_FILE_NAME)

	return &types.ManagerConfig{
		Home:       homePath,
		ConfigPath: cfgPath,
	}, nil

}

// GenerateConfigFile writes cfg to <cfg.Home>/config.toml.
// It creates the directory if needed and returns the absolute path.
func GenerateConfigFile(cfg *types.ManagerConfig) (string, error) {
	if cfg == nil {
		return "", errors.New("cfg is nil")
	}
	if cfg.Home == "" {
		return "", errors.New("cfg.Home is empty")
	}

	if cfg.ConfigPath == "" {
		return "", errors.New("cfg.COnfigPath is empty")
	}
	// Ensure directory exists
	if err := os.MkdirAll(cfg.Home, 0o755); err != nil {
		return "", err
	}

	// Marshal to TOML
	b, err := toml.Marshal(cfg)
	if err != nil {
		return "", err
	}

	// Write file
	path := cfg.ConfigPath
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return "", err
	}
	return path, nil
}
