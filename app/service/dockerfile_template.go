package service

import (
	"bytes"
	"context"
	"text/template"

	"github.com/ulysseskk/builder/app/common/log"
)

func init() {
	temp := template.New("standardDockerFile")
	temp, err := temp.Parse(DockerFileTemplate)
	if err != nil {
		panic(err)
	}
	StandardTemplate = temp
}

func GenerateStandardDockerfile(ctx context.Context, data *StandardDockerfileData) (string, error) {
	result := &bytes.Buffer{}
	err := StandardTemplate.Execute(result, data)
	if err != nil {
		log.WithContext(ctx).WithError(err).Errorf("渲染dockerfile失败")
		return "", err
	}
	return result.String(), nil
}

var StandardTemplate *template.Template

type StandardDockerfileData struct {
	RepoPath    string `json:"repo_path"`
	ServiceName string `json:"service_name"`
}

const DockerFileTemplate = `FROM golang:1.15 AS builder
# 按需安装依赖包
# RUN  apk --update --no-cache add gcc libc-dev ca-certificates
# 设置Go编译参数
ARG REPO_PATH
WORKDIR $GOPATH/src/{{.RepoPath}}
COPY . .

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.io,direct

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /main  $GOPATH/src/{{.RepoPath}}/cmd/{{.ServiceName}}/

FROM scratch

WORKDIR /app

COPY --from=builder /main .
ADD cacert.pem /etc/ssl/certs/


ENTRYPOINT ["/app/main"]`
