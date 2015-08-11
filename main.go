package main

import (
	"os"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bltconfig "github.com/mariash/bosh-load-tests/config"
	bltenv "github.com/mariash/bosh-load-tests/environment"
)

func main() {
	if len(os.Args) != 2 {
		println("Usage: blt path/to/config.json")
		os.Exit(1)
	}

	logger := boshlog.NewLogger(boshlog.LevelError)
	fs := boshsys.NewOsFileSystem(logger)
	cmdRunner := boshsys.NewExecCmdRunner(logger)

	config := bltconfig.NewConfig(fs)
	err := config.Load(os.Args[1])
	if err != nil {
		panic(err)
	}

	environmentProvider := bltenv.NewProvider(config, fs, cmdRunner)
	environment := environmentProvider.Get()
	err = environment.Setup()
	if err != nil {
		panic(err)
	}
	defer environment.Shutdown()
}
