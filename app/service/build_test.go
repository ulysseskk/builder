package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/ulysseskk/builder/app/common/connector"

	"github.com/ulysseskk/builder/app/common/config"
	"github.com/ulysseskk/builder/app/common/log"
)

func TestPrepareCodeTar(t *testing.T) {
	config.SetGlobalConfig(&config.Config{
		Github: nil,
		Log: &config.LogConfig{
			Method:   "stdout",
			FilePath: "",
			Syslog:   nil,
		},
	})
	err := log.InitLogger()
	if err != nil {
		panic(err)
	}
	tarPath, err := PrepareCodeTar(context.Background(), "house_scrapper", "9a23220ff2597dbea933fd70e43f3a3fc1669cf9", "github.com/ulysseskk/house", "/Users/konghaishuo/Documents/go/src/github.com/ulysseskk/builder/app/service/codes/9a23220ff2597dbea933fd70e43f3a3fc1669cf9")
	if err != nil {
		panic(err)
	}
	fmt.Println(tarPath)
}

func TestBuildBinaryByContainer(t *testing.T) {
	err := connector.InitDockerClient(&config.DockerConfig{Host: "tcp://192.168.50.106:2375"})
	if err != nil {
		panic(err)
	}
	resp, err := BuildBinaryByContainer(context.Background(), "/Users/konghaishuo/Documents/go/src/github.com/ulysseskk/builder/app/service/code_tars/house_scrapper/9a23220ff2597dbea933fd70e43f3a3fc1669cf9.tar", []string{"house_scrapper:v0.0.1"})
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		panic(err)
	}
}
