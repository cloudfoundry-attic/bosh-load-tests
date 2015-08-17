package dummy

import (
	"time"

	bltassets "github.com/cloudfoundry-incubator/bosh-load-tests/assets"
	bltcom "github.com/cloudfoundry-incubator/bosh-load-tests/command"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type DirectorService struct {
	directorMigrationCommand string
	directorStartCommand     string
	workerStartCommand       string
	assetsProvider           bltassets.Provider
	directorConfig           *DirectorConfig
	cmdRunner                boshsys.CmdRunner
	directorProcess          boshsys.Process
	workerProcesses          []boshsys.Process
	redisProcess             boshsys.Process
}

func NewDirectorService(
	directorMigrationCommand string,
	directorStartCommand string,
	workerStartCommand string,
	directorConfig *DirectorConfig,
	cmdRunner boshsys.CmdRunner,
	assetsProvider bltassets.Provider,
) *DirectorService {
	return &DirectorService{
		directorMigrationCommand: directorMigrationCommand,
		directorStartCommand:     directorStartCommand,
		workerStartCommand:       workerStartCommand,
		directorConfig:           directorConfig,
		cmdRunner:                cmdRunner,
		assetsProvider:           assetsProvider,
	}
}

func (s *DirectorService) Start() error {
	err := s.directorConfig.Write()
	if err != nil {
		return err
	}

	migrationCommand := bltcom.CreateCommand(s.directorMigrationCommand)
	migrationCommand.Args = append(migrationCommand.Args, "-c", s.directorConfig.DirectorConfigPath())
	_, _, _, err = s.cmdRunner.RunComplexCommand(migrationCommand)
	if err != nil {
		return bosherr.WrapError(err, "running migrations")
	}

	directorCommand := bltcom.CreateCommand(s.directorStartCommand)
	directorCommand.Args = append(directorCommand.Args, "-c", s.directorConfig.DirectorConfigPath())
	s.directorProcess, err = s.cmdRunner.RunComplexCommandAsync(directorCommand)
	if err != nil {
		return bosherr.WrapError(err, "starting director")
	}

	s.directorProcess.Wait()

	redisConfigPath := s.assetsProvider.FullPath("redis.conf")
	redisCommand := boshsys.Command{
		Name: "redis-server",
		Args: []string{redisConfigPath},
	}
	s.redisProcess, err = s.cmdRunner.RunComplexCommandAsync(redisCommand)
	if err != nil {
		return bosherr.WrapError(err, "starting redis")
	}

	s.redisProcess.Wait()

	for i := 1; i <= 3; i++ {
		workerStartCommand := bltcom.CreateCommand(s.workerStartCommand)
		workerStartCommand.Env["QUEUE"] = "*"
		workerStartCommand.Args = append(workerStartCommand.Args, "-c", s.directorConfig.WorkerConfigPath(i))

		workerProcess, err := s.cmdRunner.RunComplexCommandAsync(workerStartCommand)
		if err != nil {
			return bosherr.WrapError(err, "starting worker")
		}
		workerProcess.Wait()
		s.workerProcesses = append(s.workerProcesses, workerProcess)
	}

	return nil
}

func (s *DirectorService) Stop() {
	for _, process := range s.workerProcesses {
		process.TerminateNicely(5 * time.Second)
	}
	s.directorProcess.TerminateNicely(5 * time.Second)
	s.redisProcess.TerminateNicely(5 * time.Second)
}
