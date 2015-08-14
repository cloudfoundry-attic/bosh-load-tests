package action

import (
	"bytes"
	"path/filepath"
	"text/template"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type deploy struct {
	directorInfo   DirectorInfo
	deploymentName string
	cliRunner      *CliRunner
	fs             boshsys.FileSystem
}

type manifestData struct {
	DeploymentName string
	DirectorUUID   string
}

func NewDeploy(directorInfo DirectorInfo, deploymentName string, cliRunner *CliRunner, fs boshsys.FileSystem) *deploy {
	return &deploy{
		directorInfo:   directorInfo,
		deploymentName: deploymentName,
		cliRunner:      cliRunner,
		fs:             fs,
	}
}

func (d *deploy) Execute() error {
	err := d.cliRunner.TargetAndLogin(d.directorInfo.URL)
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
		DeploymentName: d.deploymentName,
		DirectorUUID:   d.directorInfo.UUID,
	}
	err = t.Execute(buffer, data)
	if err != nil {
		return err
	}
	err = d.fs.WriteFile(manifestPath.Name(), buffer.Bytes())
	if err != nil {
		return err
	}

	err = d.cliRunner.RunWithArgs("deployment", manifestPath.Name())
	if err != nil {
		return err
	}

	err = d.cliRunner.RunWithArgs("deploy")
	if err != nil {
		return err
	}

	return nil
}
