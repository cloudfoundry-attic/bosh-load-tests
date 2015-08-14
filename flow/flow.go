package flow

import (
	"strings"

	bltaction "github.com/cloudfoundry-incubator/bosh-load-tests/action"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
)

type actionsFlow struct {
	cliPath       string
	actionNames   []string
	actionFactory bltaction.Factory
	cmdRunner     boshsys.CmdRunner
	fs            boshsys.FileSystem
}

func NewFlow(
	cliPath string,
	actionNames []string,
	actionFactory bltaction.Factory,
	cmdRunner boshsys.CmdRunner,
	fs boshsys.FileSystem,
) *actionsFlow {
	return &actionsFlow{
		cliPath:       cliPath,
		actionNames:   actionNames,
		actionFactory: actionFactory,
		cmdRunner:     cmdRunner,
		fs:            fs,
	}
}

func (f *actionsFlow) Run() error {
	uuid, err := boshuuid.NewGenerator().Generate()
	if err != nil {
		return err
	}
	deploymentName := strings.Join([]string{"deployment", uuid}, "-")

	cliRunner := bltaction.NewCliRunner(f.cliPath, f.cmdRunner, f.fs)
	cliRunner.Configure()
	defer cliRunner.Clean()

	for _, actionName := range f.actionNames {
		action, err := f.actionFactory.Create(actionName, deploymentName, cliRunner)
		if err != nil {
			return err
		}

		err = action.Execute()
		if err != nil {
			return err
		}
	}

	return nil
}
