package server

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/naag/f1-api/internal"
	"github.com/naag/f1-api/pkg/api"
	"github.com/naag/f1-api/pkg/logging"
	"github.com/naag/f1-api/pkg/rpc"
)

type PluginServer struct {
	scenarios map[string]*scenario
	logger    hclog.Logger
	verbose   bool
}

type scenario struct {
	setupFn    api.SetupFn
	runFn      api.RunFn
	teardownFn api.TeardownFn
	t          *api.T
}

func NewServer() *PluginServer {
	p := &PluginServer{
		verbose:   false,
		scenarios: make(map[string]*scenario),
	}

	if len(os.Args) >= 2 && os.Args[1] == "-v" {
		p.WithVerboseLogging(true)
	}

	return p
}

func (p *PluginServer) WithVerboseLogging(verbose bool) *PluginServer {
	p.verbose = verbose
	return p
}

func (p *PluginServer) WithScenario(name string, setupFn api.SetupFn) *PluginServer {
	p.scenarios[name] = &scenario{
		setupFn: setupFn,
	}
	return p
}

func (p *PluginServer) Serve() {
	p.logger = logging.GetServerLogger()

	pluginMap := map[string]plugin.Plugin{
		internal.PluginName: &rpc.Plugin{Impl: p},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: internal.HandshakeConfig,
		Plugins:         pluginMap,
		Logger:          p.logger,
	})
}
