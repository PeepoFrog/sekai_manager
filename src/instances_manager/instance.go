package instancesmanager

import (
	"github.com/PeepoFrog/sekai_manager/src/cfg"
	"github.com/PeepoFrog/sekai_manager/src/types"
)

func NewInstanceManager() (*types.ManagerConfig, error) {
	return cfg.DefaultCfg()
}
