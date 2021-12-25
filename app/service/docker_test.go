package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/ulysseskk/builder/app/common/connector"

	"github.com/ulysseskk/builder/app/common/config"
)

func TestListDockerImage(t *testing.T) {
	err := connector.InitDockerClient(&config.DockerConfig{Host: "tcp://192.168.50.106:2375"})
	images, err := ListDockerImage(context.Background())
	if err != nil {
		panic(err)
	}
	for _, image := range images {
		fmt.Println(image.ID)
	}
}
