FROM golang:1.18 as builder
WORKDIR /go/src/kubefilebrowser
COPY . /go/src/kubefilebrowser

ENV PATH=$GOPATH/bin:$PATH

RUN apt-get update \
    && apt-get install musl-dev git -y

# build code
RUN go install github.com/swaggo/swag/cmd/swag@latest \
    && swag init -g cmd/server/main.go \
    && chmod +x *.sh \
    && ./00-build_lib.sh

RUN GO_VERSION=`go version|awk '{print $3" "$4}'` \
    && GIT_URL=`git remote -v|grep push|awk '{print $2}'` \
    && GIT_BRANCH=`git rev-parse --abbrev-ref HEAD` \
    && GIT_COMMIT=`git rev-parse HEAD` \
    && BUILD_TIME=`date +"%Y-%m-%d %H:%M:%S %Z"` \
    && go mod tidy \
    && CGO_ENABLED=0 go build -ldflags \
    "-w -s -X 'github.com/xmapst/kubefilebrowser.GoVersion=${GO_VERSION}' \
    -X 'github.com/xmapst/kubefilebrowser.GitUrl=${GIT_URL}' \
    -X 'github.com/xmapst/kubefilebrowser.GitBranch=${GIT_BRANCH}' \
    -X 'github.com/xmapst/kubefilebrowser.GitCommit=${GIT_COMMIT}' \
    -X 'github.com/xmapst/kubefilebrowser.BuildTime=${BUILD_TIME}'" -o main cmd/server/main.go

#1 ----------------------------
FROM alpine:latest
COPY --from=builder --chmod=0777 /go/src/kubefilebrowser/main /usr/local/bin/kubefilebrowser
RUN apk add --no-cache ca-certificates mailcap openssh jq curl busybox-extras

ENTRYPOINT ["/usr/local/bin/kubefilebrowser"]
