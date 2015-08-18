package flow

import (
	"math/rand"

	bltaction "github.com/cloudfoundry-incubator/bosh-load-tests/action"
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type ActionInfo struct {
	Name                string `json:"name"`
	DelayInMilliseconds int64  `json:"delay"`
}

type randomizer struct {
	actionFactory          bltaction.Factory
	cliRunnerFactory       bltclirunner.Factory
	state                  [][]ActionInfo
	maxDelayInMilliseconds int64
	logger                 boshlog.Logger
}

type Randomizer interface {
	Prepare(flows [][]string) error
	RunFlow(flowNumber int) error
}

func NewRandomizer(
	actionFactory bltaction.Factory,
	cliRunnerFactory bltclirunner.Factory,
	logger boshlog.Logger,
) Randomizer {
	return &randomizer{
		actionFactory:    actionFactory,
		cliRunnerFactory: cliRunnerFactory,
		state:            [][]ActionInfo{},
		maxDelayInMilliseconds: 5000,
		logger:                 logger,
	}
}

func (r *randomizer) Configure(filePath string) error {
	return nil
}

func (r *randomizer) Prepare(flows [][]string) error {
	for _, actionNames := range flows {
		actionInfos := []ActionInfo{}
		for _, actionName := range actionNames {
			actionInfos = append(actionInfos, ActionInfo{
				Name:                actionName,
				DelayInMilliseconds: rand.Int63n(r.maxDelayInMilliseconds),
			})
		}
		r.state = append(r.state, actionInfos)
	}

	r.logger.Debug("randomizer", "Generated state %#v", r.state)

	return nil
}

func (r *randomizer) RunFlow(flowNumber int) error {
	actionNames := r.state[flowNumber]
	r.logger.Debug("randomizer", "Creating flow with %#v", actionNames)

	flow := NewFlow(flowNumber, actionNames, r.actionFactory, r.cliRunnerFactory)

	return flow.Run()
}
