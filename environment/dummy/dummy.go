package dummy

import (
	"time"

	bltconfig "github.com/cloudfoundry-incubator/bosh-load-tests/config"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type dummy struct {
	workingDir      string
	database        *Database
	directorService *DirectorService
	nginxService    *NginxService
	natsService     *NatsService
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

	d.database = NewDatabase("test", d.cmdRunner)
	err = d.database.Create()
	if err != nil {
		return err
	}

	d.natsService = NewNatsService(d.config.NatsStartCommand, 65010, d.cmdRunner)
	err = d.natsService.Start()
	if err != nil {
		return err
	}

	d.nginxService = NewNginxService(d.config.NginxStartCommand, 65001, 65002, d.cmdRunner)
	err = d.nginxService.Start()
	if err != nil {
		return err
	}

	directorConfig := NewDirectorConfig(65001, d.workingDir, d.fs)
	d.directorService = NewDirectorService(
		d.config.DirectorMigrationCommand,
		d.config.DirectorStartCommand,
		d.config.WorkerStartCommand,
		directorConfig,
		d.cmdRunner,
	)

	err = d.directorService.Start()
	if err != nil {
		return err
	}

	// FIXME: wait for startup instead
	time.Sleep(5 * time.Second)

	return nil
}

func (d *dummy) Shutdown() error {
	d.nginxService.Stop()
	d.directorService.Stop()
	d.natsService.Stop()
	d.database.Drop()
	d.fs.RemoveAll(d.workingDir)

	return nil
}

func (d *dummy) DirectorURL() string {
	return "http://localhost:65002"
}
