package service

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/ulysseskk/builder/app/common/connector"
)

func ListDockerImage(ctx context.Context) ([]types.ImageSummary, error) {
	summary, err := connector.GetDockerClient().ImageList(ctx, types.ImageListOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}
	return summary, nil
}
