#docker build --network host --rm --build-arg APP_ROOT=/go/src/kubefilebrowser -t kubefilebrowser:latest -f Dockerfile .
#0 ----------------------------
FROM golang:1.18
ARG  APP_ROOT
WORKDIR ${APP_ROOT}
COPY ./ ${APP_ROOT}

ENV GO111MODULE=on
#ENV GOPROXY=https://mirrors.aliyun.com/goproxy/
ENV PATH=$GOPATH/bin:$PATH

# install upx
RUN sed -i "s/deb.debian.org/mirrors.aliyun.com/g" /etc/apt/sources.list \
  && sed -i "s/security.debian.org/mirrors.aliyun.com/g" /etc/apt/sources.list \
  && apt-get update \
  && apt-get install upx musl-dev git -y

# build code
RUN go get -u github.com/swaggo/swag/cmd/swag \
  && swag init -g cmd/server/main.go \
  && chmod +x *.sh \
  && ./00-build_lib.sh \
  #&& ./01-build_server.sh \
  && GO_VERSION=`go version|awk '{print $3" "$4}'` \
  && GIT_URL=`git remote -v|grep push|awk '{print $2}'` \
  && GIT_BRANCH=`git rev-parse --abbrev-ref HEAD` \
  && GIT_COMMIT=`git rev-parse HEAD` \
  && VERSION=`git describe --tags --abbrev=0` \
  && GIT_LATEST_TAG=`git describe --tags --abbrev=0` \
  && BUILD_TIME=`date +"%Y-%m-%d %H:%M:%S %Z"` \
  && go mod tidy \
  && CGO_ENABLED=0 GOOS=linux go build -ldflags \
  "-w -s -X 'main.Version=${VERSION}' \
  -X 'github.com/xmapst/kubefilebrowser.GoVersion=${GO_VERSION}' \
  -X 'github.com/xmapst/kubefilebrowser.GitUrl=${GIT_URL}' \
  -X 'github.com/xmapst/kubefilebrowser.GitBranch=${GIT_BRANCH}' \
  -X 'github.com/xmapst/kubefilebrowser.GitCommit=${GIT_COMMIT}' \
  -X 'github.com/xmapst/kubefilebrowser.GitLatestTag=${GIT_LATEST_TAG}' \
  -X 'github.com/xmapst/kubefilebrowser.BuildTime=${BUILD_TIME}'" -o main cmd/server/main.go \
  && strip --strip-unneeded main \
  && upx --lzma main

#1 ----------------------------
FROM alpine:latest
ARG APP_ROOT
WORKDIR /app
COPY --from=0 ${APP_ROOT}/main .
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
  && apk add --no-cache openssh jq curl busybox-extras \
  && rm -rf /var/cache/apk/*

ENTRYPOINT ["/app/main"]
