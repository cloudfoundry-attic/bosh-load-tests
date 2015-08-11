package deployment

import (
	"path/filepath"
)

type deployment struct {
	target    string
	cliRunner *CliRunner
}

func NewDeployment(target string, cliRunner *CliRunner) *deployment {
	return &deployment{
		target:    target,
		cliRunner: cliRunner,
	}
}

func (d *deployment) Deploy() error {
	err := d.cliRunner.TargetAndLogin(d.target)
	if err != nil {
		return err
	}

	cloudConfigPath, err := filepath.Abs("./assets/cloud_config.yml")
	if err != nil {
		return err
	}

	err = d.cliRunner.RunWithArgs("update", "cloud-config", cloudConfigPath)
	if err != nil {
		return err
	}

	manifestPath, err := filepath.Abs("./assets/manifest.yml")
	if err != nil {
		return err
	}

	err = d.cliRunner.RunWithArgs("deployment", manifestPath)
	if err != nil {
		return err
	}

	return nil
}
