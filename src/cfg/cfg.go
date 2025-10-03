package cfg

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

type AddressBinding struct {
	ApiAddress      string //app.toml:[api]:address 		Default: `"tcp://localhost:1317"`
	RossettaAddress string //app.toml:[rossetta]:address 	Default: `":8080"`
	GrpcAddress     string //app.toml:[grpc]:address 		Default: `"localhost:9090"`
	GrpcWebAddress  string //app.toml:[grpc-web]:address 	Default: `"localhost:9091"`

	ProxyApp                            string //config.toml:[]:proxy_app								Default: `"tcp://127.0.0.1:26658"`
	RpcLaddr                            string //config.toml:[rpc]:laddr								Default: `"tcp://127.0.0.1:26657"`
	RpcPprofLaddr                       string //config.toml:[rpc]:pprof_laddr							Default: `"localhost:6060"`
	P2PLaddr                            string //config.toml:[p2p]:laddr								Default: `"tcp://0.0.0.0:26656"`
	InstrumentationPrometheusListenAddr string //config.toml:[instrumentation]:prometheus_listen_addr	Default: `":26660"`

	Node string //client.toml:[]:node Default: `"tcp://localhost:26657"` - copy value from RpcLaddr
}

// DefaultAddressBinding returns the struct filled with defaults.
// NOTE: Node is set to match RpcLaddr by default.
func DefaultAddressBinding() AddressBinding {
	ab := AddressBinding{
		ApiAddress:                          "tcp://localhost:1317",
		RossettaAddress:                     ":8080",
		GrpcAddress:                         "localhost:9090",
		GrpcWebAddress:                      "localhost:9091",
		ProxyApp:                            "tcp://127.0.0.1:26658",
		RpcLaddr:                            "tcp://127.0.0.1:26657",
		RpcPprofLaddr:                       "localhost:6060",
		P2PLaddr:                            "tcp://0.0.0.0:26656",
		InstrumentationPrometheusListenAddr: ":26660",
	}
	ab.Node = ab.RpcLaddr
	return ab
}

// Option is a functional option for AddressBinding.
type Option func(*AddressBinding)

// NewAddressBinding returns a binding where any subset of fields can be customized.
// If Node is not explicitly set via options, it will mirror RpcLaddr after options apply.
func NewAddressBinding(opts ...Option) AddressBinding {
	ab := DefaultAddressBinding()
	for _, opt := range opts {
		opt(&ab)
	}
	// Keep Node in sync with RpcLaddr unless explicitly set by an option.
	if ab.Node == "" {
		ab.Node = ab.RpcLaddr
	}
	return ab
}

// --- Option helpers ---

func WithAPIAddress(s string) Option      { return func(a *AddressBinding) { a.ApiAddress = s } }
func WithRossettaAddress(s string) Option { return func(a *AddressBinding) { a.RossettaAddress = s } }
func WithGRPCAddress(s string) Option     { return func(a *AddressBinding) { a.GrpcAddress = s } }
func WithGRPCWebAddress(s string) Option  { return func(a *AddressBinding) { a.GrpcWebAddress = s } }

func WithProxyApp(s string) Option { return func(a *AddressBinding) { a.ProxyApp = s } }
func WithRPC(s string) Option      { return func(a *AddressBinding) { a.RpcLaddr = s } }
func WithRPCPprof(s string) Option { return func(a *AddressBinding) { a.RpcPprofLaddr = s } }
func WithP2P(s string) Option      { return func(a *AddressBinding) { a.P2PLaddr = s } }
func WithPrometheus(s string) Option {
	return func(a *AddressBinding) { a.InstrumentationPrometheusListenAddr = s }
}

// WithNode explicitly sets the client node address (won't auto-mirror RpcLaddr).
func WithNode(s string) Option { return func(a *AddressBinding) { a.Node = s } }

// WithRPCAndNode sets both RpcLaddr and Node to the same value for convenience.
func WithRPCAndNode(s string) Option {
	return func(a *AddressBinding) {
		a.RpcLaddr = s
		a.Node = s
	}
}

// PortPair holds the default and current port for a given address.
type PortPair struct {
	Default int
	Current int
}

// NamedPortPair is handy if you want a stable, ordered slice instead of a map.
type NamedPortPair struct {
	Name    string
	Default int
	Current int
}

// PortPairs returns a map of service name -> {Default, Current} port ints.
func PortPairs(ab AddressBinding) (map[string]PortPair, error) {
	def := DefaultAddressBinding()

	type item struct {
		name string
		def  string
		cur  string
	}
	items := []item{
		{"api", def.ApiAddress, ab.ApiAddress},
		{"rosetta", def.RossettaAddress, ab.RossettaAddress},
		{"grpc", def.GrpcAddress, ab.GrpcAddress},
		{"grpc_web", def.GrpcWebAddress, ab.GrpcWebAddress},

		{"proxy_app", def.ProxyApp, ab.ProxyApp},
		{"rpc", def.RpcLaddr, ab.RpcLaddr},
		{"rpc_pprof", def.RpcPprofLaddr, ab.RpcPprofLaddr},
		{"p2p", def.P2PLaddr, ab.P2PLaddr},
		{"prometheus", def.InstrumentationPrometheusListenAddr, ab.InstrumentationPrometheusListenAddr},

		// Node defaults to RpcLaddr; still report its own pair for clarity.
		{"node", def.Node, ab.Node},
	}

	out := make(map[string]PortPair, len(items))
	var firstErr error

	for _, it := range items {
		dp, errD := portOf(it.def)
		cp, errC := portOf(it.cur)
		if errD != nil && firstErr == nil {
			firstErr = fmt.Errorf("%s default addr: %w", it.name, errD)
		}
		if errC != nil && firstErr == nil {
			firstErr = fmt.Errorf("%s current addr: %w", it.name, errC)
		}
		out[it.name] = PortPair{Default: dp, Current: cp}
	}
	return out, firstErr
}

// PortPairsList returns a stable, ordered slice of {name, Default, Current}.
func PortPairsList(ab AddressBinding) ([]NamedPortPair, error) {
	m, err := PortPairs(ab)
	if err != nil {
		// still return what we could parse
	}
	order := []string{
		"api", "rosetta", "grpc", "grpc_web",
		"proxy_app", "rpc", "rpc_pprof", "p2p", "prometheus", "node",
	}
	out := make([]NamedPortPair, 0, len(order))
	for _, name := range order {
		pp := m[name]
		out = append(out, NamedPortPair{Name: name, Default: pp.Default, Current: pp.Current})
	}
	return out, err
}

// portOf extracts the numeric port from common address forms like:
// "tcp://localhost:1317", "localhost:9090", ":8080", "tcp://0.0.0.0:26656"
func portOf(addr string) (int, error) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return 0, fmt.Errorf("empty address")
	}

	// If it looks like a URL (has a scheme), use net/url first.
	if strings.Contains(addr, "://") {
		u, err := url.Parse(addr)
		if err != nil {
			return 0, err
		}
		h := u.Host
		if h == "" && u.Opaque != "" {
			// Some odd forms may end up in Opaque, try that.
			h = u.Opaque
		}
		_, p, err := splitHostPortAny(h)
		if err != nil {
			return 0, err
		}
		return parsePort(p)
	}

	// Otherwise, treat as host:port or :port
	_, p, err := splitHostPortAny(addr)
	if err != nil {
		return 0, err
	}
	return parsePort(p)
}

func splitHostPortAny(h string) (host, port string, err error) {
	// net.SplitHostPort handles IPv6 "[::1]:8080", ":8080", and "host:8080"
	host, port, err = net.SplitHostPort(h)
	if err == nil {
		return host, port, nil
	}
	// If missing brackets around IPv6, try to add them (rare in configs).
	if strings.Count(h, ":") > 1 && !strings.HasPrefix(h, "[") {
		return net.SplitHostPort("[" + h + "]")
	}
	return "", "", err
}

func parsePort(p string) (int, error) {
	if p == "" {
		return 0, fmt.Errorf("missing port")
	}
	n, err := strconv.Atoi(p)
	if err != nil || n <= 0 || n > 65535 {
		return 0, fmt.Errorf("invalid port %q", p)
	}
	return n, nil
}

func withPort(addr string, newPort int) (string, error) {
	if newPort <= 0 || newPort > 65535 {
		return "", fmt.Errorf("invalid port: %d", newPort)
	}
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return "", fmt.Errorf("empty address")
	}
	ps := strconv.Itoa(newPort)

	// URL-like ("tcp://host:port")?
	if strings.Contains(addr, "://") {
		u, err := url.Parse(addr)
		if err != nil {
			return "", err
		}
		h := u.Host
		if h == "" && u.Opaque != "" {
			// e.g. odd forms might place host in Opaque
			h = u.Opaque
		}

		host, _, err := splitHostPortLoose(h)
		if err != nil {
			// If no port present, treat entire h as a host
			host = trimIPv6Brackets(h)
		}

		u.Host = net.JoinHostPort(trimIPv6Brackets(host), ps)
		return u.String(), nil
	}

	// Bare host:port / :port
	if strings.HasPrefix(addr, ":") {
		return ":" + ps, nil
	}
	host, _, err := splitHostPortLoose(addr)
	if err != nil {
		// Assume it's a host without port; just add one.
		host = trimIPv6Brackets(addr)
	}
	return net.JoinHostPort(trimIPv6Brackets(host), ps), nil
}

// splitHostPortLoose tries net.SplitHostPort, and if there's no port, returns an error
// that callers may treat as "host without port".
func splitHostPortLoose(h string) (string, string, error) {
	host, port, err := net.SplitHostPort(h)
	if err == nil {
		return host, port, nil
	}
	// Try bracketless IPv6 convenience
	if strings.Count(h, ":") > 1 && !strings.HasPrefix(h, "[") {
		return net.SplitHostPort("[" + h + "]")
	}
	return "", "", err
}

func trimIPv6Brackets(s string) string {
	return strings.TrimSuffix(strings.TrimPrefix(s, "["), "]")
}

// Small helper to mutate a field in-place.
func updatePort(field *string, newPort int) error {
	newAddr, err := withPort(*field, newPort)
	if err != nil {
		return err
	}
	*field = newAddr
	return nil
}

// --- Setter methods on AddressBinding ---------------------------------------------

// SetAPIPort changes the port inside ApiAddress (scheme/host preserved).
func (a *AddressBinding) SetAPIPort(p int) error { return updatePort(&a.ApiAddress, p) }

// SetRossettaPort changes the port inside RossettaAddress (scheme/host preserved).
func (a *AddressBinding) SetRossettaPort(p int) error { return updatePort(&a.RossettaAddress, p) }

// SetGRPCPort changes the port inside GrpcAddress.
func (a *AddressBinding) SetGRPCPort(p int) error { return updatePort(&a.GrpcAddress, p) }

// SetGRPCWebPort changes the port inside GrpcWebAddress.
func (a *AddressBinding) SetGRPCWebPort(p int) error { return updatePort(&a.GrpcWebAddress, p) }

// SetProxyAppPort changes the port inside ProxyApp.
func (a *AddressBinding) SetProxyAppPort(p int) error { return updatePort(&a.ProxyApp, p) }

// SetRPCPort changes the port inside RpcLaddr.
// If Node currently mirrors RpcLaddr exactly, it will be updated too to keep them in sync.
func (a *AddressBinding) SetRPCPort(p int) error {
	oldRPC := a.RpcLaddr
	if err := updatePort(&a.RpcLaddr, p); err != nil {
		return err
	}
	if a.Node == oldRPC {
		// keep Node mirrored when it matched exactly before
		_ = updatePort(&a.Node, p)
	}
	return nil
}

// SetRPCPprofPort changes the port inside RpcPprofLaddr.
func (a *AddressBinding) SetRPCPprofPort(p int) error { return updatePort(&a.RpcPprofLaddr, p) }

// SetP2PPort changes the port inside P2PLaddr.
func (a *AddressBinding) SetP2PPort(p int) error { return updatePort(&a.P2PLaddr, p) }

// SetPrometheusPort changes the port inside InstrumentationPrometheusListenAddr.
func (a *AddressBinding) SetPrometheusPort(p int) error {
	return updatePort(&a.InstrumentationPrometheusListenAddr, p)
}

// SetNodePort changes the port inside Node (client.toml).
func (a *AddressBinding) SetNodePort(p int) error { return updatePort(&a.Node, p) }

// SetPort is a convenience switch by logical name.
func (a *AddressBinding) SetPort(name string, p int) error {
	switch strings.ToLower(name) {
	case "api":
		return a.SetAPIPort(p)
	case "rosetta":
		return a.SetRossettaPort(p)
	case "grpc":
		return a.SetGRPCPort(p)
	case "grpc_web", "grpc-web":
		return a.SetGRPCWebPort(p)
	case "proxy_app", "proxy-app":
		return a.SetProxyAppPort(p)
	case "rpc":
		return a.SetRPCPort(p)
	case "rpc_pprof", "rpc-pprof", "pprof":
		return a.SetRPCPprofPort(p)
	case "p2p":
		return a.SetP2PPort(p)
	case "prometheus", "metrics":
		return a.SetPrometheusPort(p)
	case "node":
		return a.SetNodePort(p)
	default:
		return fmt.Errorf("unknown port name %q", name)
	}
}

// In cfg:

// Validate checks that each address has a parseable port (1..65535).
func (a AddressBinding) Validate() error {
	var errs []string

	check := func(name, v string) {
		if _, err := portOf(v); err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", name, err))
		}
	}

	check("ApiAddress", a.ApiAddress)
	check("RossettaAddress", a.RossettaAddress)
	check("GrpcAddress", a.GrpcAddress)
	check("GrpcWebAddress", a.GrpcWebAddress)
	check("ProxyApp", a.ProxyApp)
	check("RpcLaddr", a.RpcLaddr)
	check("RpcPprofLaddr", a.RpcPprofLaddr)
	check("P2PLaddr", a.P2PLaddr)
	check("InstrumentationPrometheusListenAddr", a.InstrumentationPrometheusListenAddr)
	check("Node", a.Node)

	if len(errs) > 0 {
		return fmt.Errorf("invalid addresses:\n  - %s", strings.Join(errs, "\n  - "))
	}
	return nil
}
