package dummy

import (
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Database struct {
	name      string
	cmdRunner boshsys.CmdRunner
}

func NewDatabase(name string, cmdRunner boshsys.CmdRunner) *Database {
	return &Database{
		name:      name,
		cmdRunner: cmdRunner,
	}
}

func (d *Database) Create() error {
	d.Drop()
	_, _, _, err := d.cmdRunner.RunCommand("psql", "-U", "postgres", "-c", "create database "+d.name+";")
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) Drop() error {
	_, _, _, err := d.cmdRunner.RunCommand("psql", "-U", "postgres", "-c", "drop database "+d.name+";")
	if err != nil {
		return err
	}
	return nil
}
