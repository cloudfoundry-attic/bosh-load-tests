package main

import (
	"os"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bltconfig "github.com/mariash/bosh-load-tests/config"
	bltdep "github.com/mariash/bosh-load-tests/deployment"
	bltenv "github.com/mariash/bosh-load-tests/environment"
)

func main() {
	if len(os.Args) != 2 {
		println("Usage: blt path/to/config.json")
		os.Exit(1)
	}

	logger := boshlog.NewLogger(boshlog.LevelDebug)
	fs := boshsys.NewOsFileSystem(logger)
	cmdRunner := boshsys.NewExecCmdRunner(logger)

	config := bltconfig.NewConfig(fs)
	err := config.Load(os.Args[1])
	if err != nil {
		panic(err)
	}

	logger.Debug("main", "Setting up environment")
	environmentProvider := bltenv.NewProvider(config, fs, cmdRunner)
	environment := environmentProvider.Get()
	err = environment.Setup()
	if err != nil {
		panic(err)
	}
	defer environment.Shutdown()

	logger.Debug("main", "Starting deploy")
	cliRunner := bltdep.NewCliRunner(config.CliCmd, cmdRunner)
	deployment := bltdep.NewDeployment(environment.DirectorURL(), cliRunner)
	err = deployment.Deploy()
	if err != nil {
		panic(err)
	}
}
