package action

import (
	bltcom "github.com/cloudfoundry-incubator/bosh-load-tests/command"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type CliRunner struct {
	configPath string
	cmd        boshsys.Command
	cmdRunner  boshsys.CmdRunner
	fs         boshsys.FileSystem
}

func NewCliRunner(cliCmd string, cmdRunner boshsys.CmdRunner, fs boshsys.FileSystem) *CliRunner {
	cmd := bltcom.CreateCommand(cliCmd)

	return &CliRunner{
		cmd:       cmd,
		cmdRunner: cmdRunner,
		fs:        fs,
	}
}

func (r *CliRunner) Configure() error {
	configFile, err := r.fs.TempFile("bosh-config")
	if err != nil {
		return err
	}
	r.configPath = configFile.Name()
	return nil
}

func (r *CliRunner) Clean() error {
	if r.configPath == "" {
		return nil
	}

	return r.fs.RemoveAll(r.configPath)
}

func (r *CliRunner) TargetAndLogin(target string) error {
	err := r.RunWithArgs("target", target)
	if err != nil {
		return err
	}

	err = r.RunWithArgs("login", "admin", "admin")
	if err != nil {
		return err
	}

	return nil
}

func (r *CliRunner) RunInDirWithArgs(dir string, args ...string) error {
	cmd := r.cliCommand(args...)
	cmd.WorkingDir = dir
	_, _, _, err := r.cmdRunner.RunComplexCommand(cmd)
	if err != nil {
		return err
	}
	return nil
}

func (r *CliRunner) RunWithArgs(args ...string) error {
	_, err := r.RunWithOutput(args...)
	return err
}

func (r *CliRunner) RunWithOutput(args ...string) (string, error) {
	stdOut, _, _, err := r.cmdRunner.RunComplexCommand(r.cliCommand(args...))
	if err != nil {
		return stdOut, err
	}

	return stdOut, nil
}

func (r *CliRunner) cliCommand(args ...string) boshsys.Command {
	cmd := r.cmd
	cmd.Args = append(cmd.Args, "-n", "-c", r.configPath)
	cmd.Args = append(cmd.Args, args...)

	return cmd
}
