package dummy

import (
	"path/filepath"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type DirectorConfig struct {
	baseDir string
	fs      boshsys.FileSystem
}

func NewDirectorConfig(directorPort int, baseDir string, fs boshsys.FileSystem) *DirectorConfig {
	return &DirectorConfig{
		baseDir: baseDir,
		fs:      fs,
	}
}

func (c *DirectorConfig) Path() string {
	return filepath.Join(c.baseDir, "director.yml")
}

func (c *DirectorConfig) Write() error {
	directorTemplatePath, err := filepath.Abs("./environment/dummy/director.yml")
	if err != nil {
		return err
	}

	println("Copying to" + c.baseDir)
	return c.fs.CopyFile(directorTemplatePath, c.Path())
}
