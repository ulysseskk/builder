package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/ulysseskk/builder/app/common/constant"

	"github.com/ulysseskk/builder/app/common/log"

	"github.com/ulysseskk/builder/app/common/model/db"

	"github.com/ulysseskk/builder/app/common/config"
)

func CloneAndBuild(ctx context.Context, serviceName, version, cloneUrl, commitId, repoPath string) (*db.BuildJob, error) {
	now := time.Now()
	job := &db.BuildJob{
		ServiceName: serviceName,
		Version:     version,
		CommitId:    commitId,
		StartAt:     &now,
	}
	// 先克隆代码
	err := PullCode(ctx, serviceName, cloneUrl, config.GlobalConfig().Github.SSHPrivateKey, commitId)
	if err != nil {
		return nil, err
	}
	// 准备代码tar包
	tarPath, err := PrepareCodeTar(ctx, serviceName, commitId, repoPath, fmt.Sprintf("%s/codes/%s/%s", constant.BasePath, serviceName, commitId))
	if err != nil {
		return nil, err
	}
	resp, err := BuildBinaryByContainer(ctx, tarPath, []string{fmt.Sprintf("%s:%s", serviceName, version)})
	if err != nil {
		return nil, err
	}
	err = job.Create(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("创建job失败")
		return nil, err
	}
	job.BuildLogReader = resp.Body
	go AsyncListenForJob(ctx, job)
	return job, nil
}

func AsyncListenForJob(ctx context.Context, job *db.BuildJob) error {
	buf := make([]byte, 32*1024)
	dst := &bytes.Buffer{}
	for {
		nr, er := job.BuildLogReader.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					return ew
				}
			}
			if ew != nil {
				return ew
			}
			if nr != nw {
				return ew
			}
			job.BuildLog = fmt.Sprintf("%s%s", job.BuildLog, dst.String())
			dst.Reset()
			err := job.Update(ctx)
			if err != nil {
				log.WithContext(ctx).WithError(err).Errorf("更新job失败.")
				return err
			}
		}
		if er != nil {
			if er != io.EOF {
				return er
			}
			break
		}
	}
	now := time.Now()
	job.EndAt = &now
	err := job.Update(ctx)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("更新job失败.")
		return err
	}
	return nil
}
