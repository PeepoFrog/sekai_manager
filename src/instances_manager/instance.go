package instancesmanager

import (
	"github.com/PeepoFrog/sekai_manager/src/cfg"
	"github.com/PeepoFrog/sekai_manager/src/types"
)

type InstanceManager struct {
	*types.ManagerConfig
}

func NewInstanceManager() (*InstanceManager, error) {
	ic, err := cfg.DefaultCfg()
	if err != nil {
		return nil, err
	}
	return &InstanceManager{ManagerConfig: ic}, nil
}

func (im *InstanceManager) CreateInstance(name string) error {
	return nil
}

func (im *InstanceManager) ListInstances() (*[]types.InstanceConfig, error) {
	return nil, nil
}
