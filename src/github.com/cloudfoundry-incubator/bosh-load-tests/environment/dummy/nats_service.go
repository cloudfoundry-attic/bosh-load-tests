package dummy

import (
	"strconv"
	"time"

	bltcom "github.com/cloudfoundry-incubator/bosh-load-tests/command"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type NatsService struct {
	natsStartCommand string
	natsPort         int
	cmdRunner        boshsys.CmdRunner
	process          boshsys.Process
}

func NewNatsService(
	natsStartCommand string,
	natsPort int,
	cmdRunner boshsys.CmdRunner,
) *NatsService {
	return &NatsService{
		natsStartCommand: natsStartCommand,
		natsPort:         natsPort,
		cmdRunner:        cmdRunner,
	}
}

func (s *NatsService) Start() error {
	var err error
	natsStartCommand := bltcom.CreateCommand(s.natsStartCommand)
	natsStartCommand.Args = append(natsStartCommand.Args, "-p", strconv.Itoa(s.natsPort), "-T")
	s.process, err = s.cmdRunner.RunComplexCommandAsync(natsStartCommand)
	if err != nil {
		return bosherr.WrapError(err, "starting nats")
	}

	s.process.Wait()

	return nil
}

func (s *NatsService) Stop() {
	s.process.TerminateNicely(5 * time.Second)
}
