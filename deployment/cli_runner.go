package deployment

import (
	bltcom "github.com/cloudfoundry-incubator/bosh-load-tests/command"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type CliRunner struct {
	cmd       boshsys.Command
	cmdRunner boshsys.CmdRunner
}

func NewCliRunner(cliCmd string, cmdRunner boshsys.CmdRunner) *CliRunner {
	cmd := bltcom.CreateCommand(cliCmd)

	return &CliRunner{
		cmd:       cmd,
		cmdRunner: cmdRunner,
	}
}

func (r *CliRunner) TargetAndLogin(target string) error {
	err := r.RunWithArgs("-n", "target", target)
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
	cmd := r.cmd
	cmd.Args = append(cmd.Args, args...)
	cmd.WorkingDir = dir
	_, _, _, err := r.cmdRunner.RunComplexCommand(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (r *CliRunner) RunWithOutput(args ...string) (string, error) {
	cmd := r.cmd
	cmd.Args = append(cmd.Args, args...)
	stdOut, _, _, err := r.cmdRunner.RunComplexCommand(cmd)
	if err != nil {
		return stdOut, err
	}

	return stdOut, nil
}

func (r *CliRunner) RunWithArgs(args ...string) error {
	_, err := r.RunWithOutput(args...)
	return err
}
