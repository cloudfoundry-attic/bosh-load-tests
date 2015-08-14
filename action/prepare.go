package action

import (
	"path/filepath"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type prepare struct {
	directorInfo DirectorInfo
	cliRunner    *CliRunner
	fs           boshsys.FileSystem
}

func NewPrepare(directorInfo DirectorInfo, cliRunner *CliRunner, fs boshsys.FileSystem) *prepare {
	return &prepare{
		directorInfo: directorInfo,
		cliRunner:    cliRunner,
		fs:           fs,
	}
}

func (p *prepare) Execute() error {
	err := p.cliRunner.TargetAndLogin(p.directorInfo.URL)
	if err != nil {
		return err
	}

	cloudConfigPath, err := filepath.Abs("./assets/cloud_config.yml")
	if err != nil {
		return err
	}

	err = p.cliRunner.RunWithArgs("update", "cloud-config", cloudConfigPath)
	if err != nil {
		return err
	}

	stemcellPath, err := filepath.Abs("./assets/stemcell.tgz")
	if err != nil {
		return err
	}

	err = p.cliRunner.RunWithArgs("upload", "stemcell", stemcellPath)
	if err != nil {
		return err
	}

	releaseDir, err := p.fs.TempDir("release-test")
	if err != nil {
		return err
	}
	defer p.fs.RemoveAll(releaseDir)

	releaseSrcPath, err := filepath.Abs("./assets/release")
	if err != nil {
		return err
	}

	err = p.fs.CopyDir(releaseSrcPath, releaseDir)
	if err != nil {
		return err
	}

	err = p.cliRunner.RunInDirWithArgs(releaseDir, "create", "release", "--force")
	if err != nil {
		return err
	}

	err = p.cliRunner.RunInDirWithArgs(releaseDir, "upload", "release")
	if err != nil {
		return err
	}
	return nil
}
