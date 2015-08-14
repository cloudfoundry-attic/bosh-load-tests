package action

import (
	"errors"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Factory interface {
	Create(name string, deploymentName string, cliRunner *CliRunner) (Action, error)
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

func (f *factory) Create(name string, deploymentName string, cliRunner *CliRunner) (Action, error) {
	switch name {
	case "prepare":
		return NewPrepare(f.directorInfo, cliRunner, f.fs), nil
	case "deploy":
		return NewDeploy(f.directorInfo, deploymentName, cliRunner, f.fs), nil
	}

	return nil, errors.New("unknown action")
}
