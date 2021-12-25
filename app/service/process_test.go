package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ulysseskk/builder/app/common/constant"

	"github.com/ulysseskk/builder/app/common/connector"

	"github.com/ulysseskk/builder/app/common/config"
)

func TestCloneAndBuild(t *testing.T) {
	constant.BasePath = "/Users/konghaishuo/Documents/go/src/github.com/ulysseskk/builder/app/service"
	config.SetGlobalConfig(&config.Config{
		Github: &config.GithubConfig{
			Token:         "",
			UseProxy:      false,
			Proxy:         nil,
			SSHPrivateKey: sshPrivateKey,
		},
		Log: &config.LogConfig{
			Method:   "stdout",
			FilePath: "",
			Syslog:   nil,
		},
		Docker: &config.DockerConfig{Host: "tcp://192.168.50.106:2375"},
		Mysql: &config.MysqlConfig{
			Host:     "192.168.50.106",
			Port:     3306,
			DbName:   "cicd",
			User:     "root",
			Password: "Khs19940718!",
		},
	})
	connector.InitMysql()
	err := connector.InitDockerClient(config.GlobalConfig().Docker)
	if err != nil {
		panic(err)
	}
	job, err := CloneAndBuild(context.Background(), "house_scrapper", "v0.0.2", "git@github.com:ulysseskk/house.git", "9a23220ff2597dbea933fd70e43f3a3fc1669cf9", "github.com/ulysseskk/house")
	if err != nil {
		panic(err)
	}
	for job.EndAt == nil {
		time.Sleep(1 * time.Second)
		fmt.Println(job.BuildLog)
	}

}
