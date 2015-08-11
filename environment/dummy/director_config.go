package dummy

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type directorConfigOptions struct {
	Port int
}

type DirectorConfig struct {
	directorPort int
	numWorkers   int
	baseDir      string
	fs           boshsys.FileSystem
}

func NewDirectorConfig(directorPort int, baseDir string, fs boshsys.FileSystem) *DirectorConfig {
	return &DirectorConfig{
		directorPort: directorPort,
		numWorkers:   3,
		baseDir:      baseDir,
		fs:           fs,
	}
}

func (c *DirectorConfig) DirectorConfigPath() string {
	return filepath.Join(c.baseDir, "director.yml")
}

func (c *DirectorConfig) WorkerConfigPath(index int) string {
	return filepath.Join(c.baseDir, fmt.Sprintf("worker-%d.yml", index))
}

func (c *DirectorConfig) Write() error {
	directorTemplatePath, err := filepath.Abs("./environment/dummy/director.yml")
	if err != nil {
		return err
	}

	t := template.Must(template.ParseFiles(directorTemplatePath))
	err = c.saveConfig(c.directorPort, c.DirectorConfigPath(), t)
	if err != nil {
		return err
	}

	for i := 1; i <= c.numWorkers; i++ {
		err = c.saveConfig(c.directorPort+i, c.WorkerConfigPath(i), t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *DirectorConfig) saveConfig(port int, path string, t *template.Template) error {
	buffer := bytes.NewBuffer([]byte{})
	options := directorConfigOptions{
		Port: port,
	}
	err := t.Execute(buffer, options)
	if err != nil {
		return err
	}
	err = c.fs.WriteFile(path, buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}
