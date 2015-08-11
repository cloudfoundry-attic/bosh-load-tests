package dummy

import (
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	bltconfig "github.com/mariash/bosh-load-tests/config"
)

type dummy struct {
	workingDir      string
	database        *Database
	directorService *DirectorService
	config          *bltconfig.Config
	fs              boshsys.FileSystem
	cmdRunner       boshsys.CmdRunner
}

func NewDummy(config *bltconfig.Config, fs boshsys.FileSystem, cmdRunner boshsys.CmdRunner) *dummy {
	return &dummy{
		config:    config,
		fs:        fs,
		cmdRunner: cmdRunner,
	}
}

func (d *dummy) Setup() error {
	var err error
	d.workingDir, err = d.fs.TempDir("dummy-working-dir")
	if err != nil {
		return err
	}

	directorConfig := NewDirectorConfig(65001, d.workingDir, d.fs)
	err = directorConfig.Write()
	if err != nil {
		return err
	}

	d.database = NewDatabase("test", d.cmdRunner)
	err = d.database.Create()
	if err != nil {
		return err
	}

	d.directorService = NewDirectorService(
		d.config.DirectorMigrationCommand,
		d.config.DirectorStartCommand,
		d.config.WorkerStartCommand,
		directorConfig.Path(),
		d.cmdRunner,
	)

	return d.directorService.Start()
	// migrate database
	// start redis
	// start director
	// wait for director to start
	// start workers
	// wait for workers to start

	// start nats

	// return nil
}

func (d *dummy) Shutdown() error {
	d.database.Drop()
	d.fs.RemoveAll(d.workingDir)

	return nil
}

func (d *dummy) DirectorURL() string {
	return "http://localhost:65001"
}
