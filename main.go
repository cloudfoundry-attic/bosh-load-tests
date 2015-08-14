package main

import (
	"os"
	"os/signal"
	"syscall"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bltaction "github.com/cloudfoundry-incubator/bosh-load-tests/action"
	bltconfig "github.com/cloudfoundry-incubator/bosh-load-tests/config"
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

	cliRunner := bltaction.NewCliRunner(config.CliCmd, cmdRunner, fs)
	cliRunner.Configure()
	defer cliRunner.Clean()

	directorInfo, err := bltaction.NewDirectorInfo(environment.DirectorURL(), cliRunner)
	if err != nil {
		panic(err)
	}

	actionFactory := bltaction.NewFactory(directorInfo, fs)

	prepareAction, _ := actionFactory.Create("prepare", "", cliRunner)
	err = prepareAction.Execute()
	if err != nil {
		panic(err)
	}

	deployAction, _ := actionFactory.Create("deploy", "my-deploy", cliRunner)
	err = deployAction.Execute()
	if err != nil {
		panic(err)
	}

	// doneCh := make(chan error)

	// for i := 0; i < config.NumberOfDeployments; i++ {
	// 	go func(i int) {
	// 		flow := bltflow.NewFlow([]string{"deploy"}, actionFactory)
	// 		doneCh <- flow.Run()
	// 	}(i)
	// }

	// for i := 0; i < config.NumberOfDeployments; i++ {
	// 	<-doneCh
	// }

	println("Done!")
}
