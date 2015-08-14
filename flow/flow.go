package flow

import (
	"strings"

	bltaction "github.com/cloudfoundry-incubator/bosh-load-tests/action"
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
)

type actionsFlow struct {
	flowNumber       int
	actionNames      []string
	actionFactory    bltaction.Factory
	cliRunnerFactory bltclirunner.Factory
}

func NewFlow(
	flowNumber int,
	actionNames []string,
	actionFactory bltaction.Factory,
	cliRunnerFactory bltclirunner.Factory,
) *actionsFlow {
	return &actionsFlow{
		flowNumber:       flowNumber,
		actionNames:      actionNames,
		actionFactory:    actionFactory,
		cliRunnerFactory: cliRunnerFactory,
	}
}

func (f *actionsFlow) Run() error {
	uuid, err := boshuuid.NewGenerator().Generate()
	if err != nil {
		return err
	}
	deploymentName := strings.Join([]string{"deployment", uuid}, "-")

	cliRunner := f.cliRunnerFactory.Create()
	cliRunner.Configure()
	defer cliRunner.Clean()

	for _, actionName := range f.actionNames {
		action, err := f.actionFactory.Create(actionName, f.flowNumber, deploymentName, cliRunner)
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
