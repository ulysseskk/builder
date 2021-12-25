package service

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ulysseskk/builder/app/common/log"

	"github.com/docker/docker/api/types"

	"github.com/ulysseskk/builder/app/common/connector"
)

func BuildBinaryByContainer(ctx context.Context, tarPath string, tags []string) (*types.ImageBuildResponse, error) {
	codeTars, err := os.Open(tarPath)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("打开tar文件失败")
		return nil, err
	}
	resp, err := connector.GetDockerClient().ImageBuild(ctx, codeTars, types.ImageBuildOptions{
		Tags:       tags,
		Dockerfile: "DockerfileAuto",
	})
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("Docker构建镜像请求失败")
		return nil, err
	}
	return &resp, nil
}

func PrepareCodeTar(ctx context.Context, serviceName, commitId, repoPath, codePath string) (string, error) {
	// 先生成dockerfile
	dockerFileContent, err := GenerateStandardDockerfile(ctx, &StandardDockerfileData{
		RepoPath:    repoPath,
		ServiceName: serviceName,
	})
	if err != nil {
		return "", err
	}
	// dockerfile写入code path
	file, err := os.Create(fmt.Sprintf("%s/DockerfileAuto", codePath))
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("创建dockerfile失败")
		return "", err
	}
	_, err = file.WriteString(dockerFileContent)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("写入dockerfile失败")
		return "", err
	}
	codePathFile, err := os.Open(codePath)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("打开code path失败")
		return "", err
	}
	err = os.MkdirAll(fmt.Sprintf("code_tars/%s/", serviceName), fs.ModePerm)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("创建代码tar包目录失败")
		return "", err
	}
	tarFile, err := os.Create(fmt.Sprintf("code_tars/%s/%s.tar", serviceName, commitId))
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("创建代码tar包失败")
		return "", err
	}
	tarWriter := tar.NewWriter(tarFile)
	defer tarWriter.Close()
	err = compress(codePathFile, "", tarWriter, true)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("打包代码tar包失败")
		return "", err
	}
	absPath, err := filepath.Abs(tarFile.Name())
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("获取tar绝对路径失败")
		return "", err
	}
	return absPath, nil
}

func compress(file *os.File, prefix string, tw *tar.Writer, withoutDirPrefix bool) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		if !withoutDirPrefix {
			prefix = prefix + "/" + info.Name()
		}
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, tw, false)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := tar.FileInfoHeader(info, "")
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(tw, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
