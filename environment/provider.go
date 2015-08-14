package environment

import (
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bltconfig "github.com/cloudfoundry-incubator/bosh-load-tests/config"
	bltdummy "github.com/cloudfoundry-incubator/bosh-load-tests/environment/dummy"
)

type provider struct {
	config       *bltconfig.Config
	environments map[string]Environment
}

func NewProvider(config *bltconfig.Config, fs boshsys.FileSystem, cmdRunner boshsys.CmdRunner) *provider {
	return &provider{
		config: config,
		environments: map[string]Environment{
			"dummy": bltdummy.NewDummy(config, fs, cmdRunner),
		},
	}
}

func (p *provider) Get() Environment {
	return p.environments[p.config.Environment]
}
