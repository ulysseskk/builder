package connector

import (
	"github.com/ulysseskk/builder/app/common/config"
)

func Init(conf *config.Config) error {
	initGithubConnector(conf.Github)
	InitMysql()
	err := InitDockerClient(conf.Docker)
	if err != nil {
		return err
	}
	return nil
}
