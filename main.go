package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bltconfig "github.com/cloudfoundry-incubator/bosh-load-tests/config"
	bltdep "github.com/cloudfoundry-incubator/bosh-load-tests/deployment"
	bltenv "github.com/cloudfoundry-incubator/bosh-load-tests/environment"
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
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		environment.Shutdown()
		os.Exit(1)
	}()

	logger.Debug("main", "Starting deploy")
	cliRunner := bltdep.NewCliRunner(config.CliCmd, cmdRunner)
	deployment := bltdep.NewDeployment(environment.DirectorURL(), cliRunner, fs)
	deployment.Prepare()

	doneCh := make(chan error)

	for i := 0; i < config.NumberOfDeployments; i++ {
		go func(i int) {
			doneCh <- deployment.Deploy(fmt.Sprintf("my-deploy-%d", i))
		}(i)
	}

	for i := 0; i < config.NumberOfDeployments; i++ {
		<-doneCh
	}

	println("Done!")
}
