package action

import (
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	bltassets "github.com/cloudfoundry-incubator/bosh-load-tests/assets"
)

type uploadCloudConfig struct {
	cliRunner      bltclirunner.Runner
	assetsProvider bltassets.Provider
}

func NewUploadCloudConfig(
	cliRunner bltclirunner.Runner,
	assetsProvider bltassets.Provider,
) *uploadCloudConfig {
	return &uploadCloudConfig{
		cliRunner:      cliRunner,
		assetsProvider: assetsProvider,
	}
}

func (p *uploadCloudConfig) Execute() error {
	cloudConfigPath, err := p.assetsProvider.FullPath("cloud_config.yml")
	if err != nil {
		return err
	}

	err = p.cliRunner.RunWithArgs("update", "cloud-config", cloudConfigPath)
	if err != nil {
		return err
	}

	return nil
}
