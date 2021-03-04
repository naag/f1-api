package server

import (
	"errors"

	"github.com/hashicorp/go-hclog"
	"github.com/naag/f1-api/pkg/api"
)

func (p *PluginServer) GetScenarios() ([]string, error) {
	var result []string

	for name := range p.scenarios {
		result = append(result, name)
	}

	return result, nil
}

func (p *PluginServer) SetupScenario(name string) error {
	s := p.getScenarioByName(name)
	p.logger = getLogger(name, p.verbose)
	s.t = newT(name, p.logger)
	s.runFn, s.teardownFn = s.setupFn(s.t)

	if s.t.HasFailed() {
		return errors.New("setup scenario failed")
	}

	return nil
}

func (p *PluginServer) RunScenarioIteration(name string) error {
	s := p.getScenarioByName(name)
	s.t = newT(name, p.logger)

	s.runFn(s.t)

	if s.t.HasFailed() {
		return errors.New("iteration failed")
	}

	return nil
}

func (p *PluginServer) StopScenario(name string) error {
	s := p.getScenarioByName(name)
	s.t = newT(name, p.logger)

	s.teardownFn(s.t)

	if s.t.HasFailed() {
		return errors.New("stop scenario failed")
	}

	return nil
}

func (p *PluginServer) getScenarioByName(name string) *scenario {
	return p.scenarios[name]
}

func newT(name string, log hclog.Logger) *api.T {
	return api.NewT("0", "0", name, log)
}
