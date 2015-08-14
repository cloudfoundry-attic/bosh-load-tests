package flow

// import (
// 	bltaction "github.com/cloudfoundry-incubator/bosh-load-tests/action"
// 	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
// )

// type deploymentFlow struct {
// 	actionNames   []string
// 	actionFactory bltaction.Factory
// }

// func NewFlow(actionNames []string, actionFactory bltaction.Factory) *deploymentFlow {
// 	return &deploymentFlow{
// 		actionNames:   actionsNames,
// 		actionFactory: actionFactory,
// 	}
// }

// func (f *deploymentFlow) Run() error {
// 	deploymentName, err := boshuuid.NewGenerator().Generate()
// 	if err != nil {
// 		return err
// 	}

// 	for _, actionName := range f.actionNames {
// 		action, err := f.actionFactory.Create(actionName)
// 		if err != nil {
// 			return err
// 		}

// 		err = action.Execute(deploymentName)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
