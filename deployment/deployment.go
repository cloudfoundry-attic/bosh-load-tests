package deployment

import (
	"bytes"
	"path/filepath"
	"text/template"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type deployment struct {
	target       string
	cliRunner    *CliRunner
	directorUUID string
	fs           boshsys.FileSystem
}

func NewDeployment(target string, cliRunner *CliRunner, fs boshsys.FileSystem) *deployment {
	return &deployment{
		target:    target,
		cliRunner: cliRunner,
		fs:        fs,
	}
}

type manifestData struct {
	DeploymentName string
	DirectorUUID   string
}

func (d *deployment) Prepare() error {
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

	stemcellPath, err := filepath.Abs("./assets/stemcell.tgz")
	if err != nil {
		return err
	}

	err = d.cliRunner.RunWithArgs("upload", "stemcell", stemcellPath)
	if err != nil {
		return err
	}

	releaseDir, err := d.fs.TempDir("release-test")
	if err != nil {
		return err
	}
	defer d.fs.RemoveAll(releaseDir)

	releaseSrcPath, err := filepath.Abs("./assets/release")
	if err != nil {
		return err
	}

	err = d.fs.CopyDir(releaseSrcPath, releaseDir)
	if err != nil {
		return err
	}

	err = d.cliRunner.RunInDirWithArgs(releaseDir, "create", "release", "--force")
	if err != nil {
		return err
	}

	err = d.cliRunner.RunInDirWithArgs(releaseDir, "upload", "release")
	if err != nil {
		return err
	}

	d.directorUUID, err = d.cliRunner.RunWithOutput("status", "--uuid")
	if err != nil {
		return err
	}

	return nil
}

func (d *deployment) Deploy(name string) error {
	configPath, err := d.fs.TempFile("bosh-config")
	if err != nil {
		return err
	}
	defer d.fs.RemoveAll(configPath.Name())

	err = d.cliRunner.RunWithArgs("-n", "-c", configPath.Name(), "target", d.target)
	if err != nil {
		return err
	}

	err = d.cliRunner.RunWithArgs("-c", configPath.Name(), "login", "admin", "admin")
	if err != nil {
		return err
	}

	manifestTemplatePath, err := filepath.Abs("./assets/manifest.yml")
	if err != nil {
		return err
	}

	manifestPath, err := d.fs.TempFile("manifest-test")
	if err != nil {
		return err
	}
	defer d.fs.RemoveAll(manifestPath.Name())

	t := template.Must(template.ParseFiles(manifestTemplatePath))
	buffer := bytes.NewBuffer([]byte{})
	data := manifestData{
		DeploymentName: name,
		DirectorUUID:   d.directorUUID,
	}
	err = t.Execute(buffer, data)
	if err != nil {
		return err
	}
	err = d.fs.WriteFile(manifestPath.Name(), buffer.Bytes())
	if err != nil {
		return err
	}

	err = d.cliRunner.RunWithArgs("-c", configPath.Name(), "deployment", manifestPath.Name())
	if err != nil {
		return err
	}

	err = d.cliRunner.RunWithArgs("-c", configPath.Name(), "-n", "deploy")
	if err != nil {
		return err
	}

	return nil
}
