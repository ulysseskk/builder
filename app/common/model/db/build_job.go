package db

import (
	"context"
	"io"
	"time"

	"github.com/ulysseskk/builder/app/common/connector"
)

type BuildJob struct {
	Id             uint64        `json:"id"`
	ServiceName    string        `json:"service_name"`
	Version        string        `json:"version"`
	CommitId       string        `json:"commit_id"`
	ContainerSha   string        `json:"container_sha"`
	Status         string        `json:"status"`
	BuildLog       string        `json:"build_log"`
	BuildLogReader io.ReadCloser `json:"-" gorm:"-"`
	StartAt        *time.Time    `json:"start_at"`
	EndAt          *time.Time    `json:"end_at"`
}

func (b *BuildJob) Create(ctx context.Context) error {
	return connector.GetMysqlConnector(ctx).Create(b).Error
}

func (b *BuildJob) Update(ctx context.Context) error {
	return connector.GetMysqlConnector(ctx).Save(b).Error
}

func (BuildJob) TableName() string {
	return "build_job"
}
