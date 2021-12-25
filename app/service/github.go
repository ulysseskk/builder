package service

import (
	"context"
	"fmt"
	"os"

	"github.com/ulysseskk/builder/app/common/log"
	"github.com/ulysseskk/builder/app/common/util"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

func PullCode(ctx context.Context, serviceName, cloneUrl string, sshKey string, commitId string) error {
	// 先尝试清理代码目录
	codePath := fmt.Sprintf("codes/%s/%s", serviceName, commitId)
	if util.Exist(codePath) {
		err := os.RemoveAll(codePath)
		if err != nil {
			log.WithContext(ctx).WithError(err).Errorf("清理代码目录失败")
			return err
		}
	}
	publicKeys, err := ssh.NewPublicKeys("git", []byte(sshKey), "")
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("解析ssh private key失败")
		return err
	}
	repo, err := git.PlainClone(fmt.Sprintf("codes/%s/%s", serviceName, commitId), false, &git.CloneOptions{
		URL:  cloneUrl,
		Auth: publicKeys,
	})
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("拉取代码失败")
		return err
	}
	workTree, err := repo.Worktree()
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("获取work tree失败")
		return err
	}
	err = workTree.Checkout(&git.CheckoutOptions{
		Hash:  plumbing.NewHash(commitId),
		Force: true,
	})
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("checkout指定commit失败")
		return err
	}
	return nil
}
