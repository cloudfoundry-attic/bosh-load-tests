package action

import (
	"errors"

	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Factory interface {
	Create(name string, flowNumber int, deploymentName string, cliRunner bltclirunner.Runner) (Action, error)
}

type factory struct {
	directorInfo DirectorInfo
	fs           boshsys.FileSystem
}

func NewFactory(
	directorInfo DirectorInfo,
	fs boshsys.FileSystem,
) *factory {
	return &factory{
		directorInfo: directorInfo,
		fs:           fs,
	}
}

func (f *factory) Create(
	name string,
	flowNumber int,
	deploymentName string,
	cliRunner bltclirunner.Runner,
) (Action, error) {
	switch name {
	case "prepare":
		return NewPrepare(f.directorInfo, cliRunner, f.fs), nil
	case "deploy_with_dynamic":
		return NewDeployWithDynamic(f.directorInfo, deploymentName, cliRunner, f.fs), nil
	case "deploy_with_static":
		return NewDeployWithStatic(f.directorInfo, flowNumber, deploymentName, cliRunner, f.fs), nil
	case "recreate":
		return NewRecreate(f.directorInfo, deploymentName, cliRunner, f.fs), nil
	case "stop_hard":
		return NewStopHard(f.directorInfo, deploymentName, cliRunner, f.fs), nil
	}

	return nil, errors.New("unknown action")
}
