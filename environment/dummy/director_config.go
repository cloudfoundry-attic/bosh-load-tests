package dummy

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type DirectorOptions struct {
	Port         int
	DatabaseName string
}

type DirectorConfig struct {
	options    DirectorOptions
	numWorkers int
	baseDir    string
	fs         boshsys.FileSystem
}

func NewDirectorConfig(options DirectorOptions, baseDir string, fs boshsys.FileSystem) *DirectorConfig {
	return &DirectorConfig{
		options:    options,
		numWorkers: 3,
		baseDir:    baseDir,
		fs:         fs,
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
	err = c.saveConfig(c.options.Port, c.DirectorConfigPath(), t)
	if err != nil {
		return err
	}

	for i := 1; i <= c.numWorkers; i++ {
		err = c.saveConfig(c.options.Port+i, c.WorkerConfigPath(i), t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *DirectorConfig) saveConfig(port int, path string, t *template.Template) error {
	buffer := bytes.NewBuffer([]byte{})
	context := c.options
	context.Port = port
	err := t.Execute(buffer, context)
	if err != nil {
		return err
	}
	err = c.fs.WriteFile(path, buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}
