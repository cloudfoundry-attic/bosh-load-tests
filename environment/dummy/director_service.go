package dummy

import (
	"fmt"
	"strings"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type DirectorService struct {
	directorMigrationCommand string
	directorStartCommand     string
	workerStartCommand       string
	configPath               string
	cmdRunner                boshsys.CmdRunner
	directorProcess          boshsys.Process
	workerProcess            boshsys.Process
}

func NewDirectorService(
	directorMigrationCommand string,
	directorStartCommand string,
	workerStartCommand string,
	configPath string,
	cmdRunner boshsys.CmdRunner,
) *DirectorService {
	return &DirectorService{
		directorMigrationCommand: directorMigrationCommand,
		directorStartCommand:     directorStartCommand,
		workerStartCommand:       workerStartCommand,
		configPath:               configPath,
		cmdRunner:                cmdRunner,
	}
}

func (s *DirectorService) Start() error {
	_, _, _, err := s.cmdRunner.RunComplexCommand(s.createCommand(s.directorMigrationCommand))
	if err != nil {
		return err
	}

	s.directorProcess, err = s.cmdRunner.RunComplexCommandAsync(s.createCommand(s.directorStartCommand))
	if err != nil {
		return err
	}

	// wait for director

	s.workerProcess, err = s.cmdRunner.RunComplexCommandAsync(s.createCommand(s.workerStartCommand))
	if err != nil {
		return err
	}

	// wait for worker

	return nil
}

func (s *DirectorService) Stop() {
	s.directorProcess.TerminateNicely(5 * time.Second)
	s.workerProcess.TerminateNicely(5 * time.Second)
}

func (s *DirectorService) createCommand(command string) boshsys.Command {
	cmdParts := strings.Split(command, " ")
	args := []string{}
	env := map[string]string{}
	var name string

	for i := 0; i < len(cmdParts); i++ {
		if strings.Contains(cmdParts[i], "=") {
			envPair := strings.Split(cmdParts[i], "=")
			env[envPair[0]] = envPair[1]
			continue
		}

		if name == "" {
			name = cmdParts[i]
			continue
		}

		args = append(args, cmdParts[i])
	}

	args = append(args, "-c", s.configPath)

	fmt.Printf("%#v", env)

	return boshsys.Command{
		Name: name,
		Args: args,
		Env:  env,
	}
}
