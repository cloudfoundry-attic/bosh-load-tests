package dummy

import (
	"path/filepath"
	"time"

	bltcom "github.com/cloudfoundry-incubator/bosh-load-tests/command"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type NginxService struct {
	nginxStartCommand string
	directorPort      int
	nginxPort         int
	cmdRunner         boshsys.CmdRunner
	process           boshsys.Process
}

func NewNginxService(
	nginxStartCommand string,
	directorPort int,
	nginxPort int,
	cmdRunner boshsys.CmdRunner,
) *NginxService {
	return &NginxService{
		nginxStartCommand: nginxStartCommand,
		directorPort:      directorPort,
		nginxPort:         nginxPort,
		cmdRunner:         cmdRunner,
	}
}

func (s *NginxService) Start() error {
	nginxStartCommand := bltcom.CreateCommand(s.nginxStartCommand)
	configPath, err := filepath.Abs("./environment/dummy/nginx.yml")
	if err != nil {
		return bosherr.WrapError(err, "getting path to nginx config")
	}

	nginxStartCommand.Args = append(nginxStartCommand.Args, "-c", configPath)
	s.process, err = s.cmdRunner.RunComplexCommandAsync(nginxStartCommand)
	if err != nil {
		return bosherr.WrapError(err, "starting nginx")
	}

	s.process.Wait()

	return nil
}

func (s *NginxService) Stop() {
	s.process.TerminateNicely(5 * time.Second)
}
