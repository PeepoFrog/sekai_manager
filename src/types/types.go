package types

// InstanceConfig describes one managed instance.
// TOML will render this as an array of tables: [[instances]]
type InstanceConfig struct {
	Name          string `toml:"name"`
	Home          string `toml:"home"`
	PortRange     int    `toml:"port_range"`
	SekaidVersion string `toml:"sekaid_version"`
}

// ManagerConfig is the root of the config file.
type ManagerConfig struct {
	Home       string           `toml:"home"`
	ConfigPath string           `toml:"config_path"`
	Instances  []InstanceConfig `toml:"instances,omitempty"`
}
