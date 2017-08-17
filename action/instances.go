package action

import (
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
)

type InstancesInfo struct {
	directorInfo   DirectorInfo
	deploymentName string
	cliRunner      bltclirunner.Runner
}

type Instance struct {
}

func NewInstances(directorInfo DirectorInfo, deploymentName string, cliRunner bltclirunner.Runner) *InstancesInfo {
	return &InstancesInfo{
		directorInfo:   directorInfo,
		deploymentName: deploymentName,
		cliRunner:      cliRunner,
	}
}

func (i *InstancesInfo) GetInstances() ([]Instance, error) {
	return []Instance{}, nil
}
