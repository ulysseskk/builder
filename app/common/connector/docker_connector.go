package connector

import (
	"github.com/docker/docker/client"
	"github.com/ulysseskk/builder/app/common/config"
)

var DockerClient *client.Client

func InitDockerClient(conf *config.DockerConfig) error {
	dockerClient, err := client.NewClientWithOpts(client.WithHost(conf.Host))
	if err != nil {
		return err
	}
	DockerClient = dockerClient
	return nil
}

func GetDockerClient() *client.Client {
	return DockerClient
}
