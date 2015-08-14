package action

import (
	"bytes"
	"net"
	"path/filepath"
	"text/template"

	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type deployWithStatic struct {
	directorInfo   DirectorInfo
	flowNumber     int
	deploymentName string
	cliRunner      bltclirunner.Runner
	fs             boshsys.FileSystem
}

func NewDeployWithStatic(
	directorInfo DirectorInfo,
	flowNumber int,
	deploymentName string,
	cliRunner bltclirunner.Runner,
	fs boshsys.FileSystem,
) *deployWithStatic {
	return &deployWithStatic{
		directorInfo:   directorInfo,
		flowNumber:     flowNumber,
		deploymentName: deploymentName,
		cliRunner:      cliRunner,
		fs:             fs,
	}
}

func (d *deployWithStatic) Execute() error {
	err := d.cliRunner.TargetAndLogin(d.directorInfo.URL)
	if err != nil {
		return err
	}

	manifestTemplatePath, err := filepath.Abs("./assets/manifest_with_static.yml")
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
		StaticIP:       d.getNextIP(),
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

func (d *deployWithStatic) getNextIP() string {
	ip := net.ParseIP("192.168.1.10")
	b := ip.To4()
	b[3] = b[3] + byte(d.flowNumber)
	return net.IPv4(b[0], b[1], b[2], b[3]).String()
}
