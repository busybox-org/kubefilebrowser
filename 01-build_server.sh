#!/usr/bin/env bash
export PATH=$PATH:$GOPATH/bin
#go get -u github.com/swaggo/swag/cmd/swag
swag init -g cmd/server/main.go
# shellcheck disable=SC2181
if [ "$?" != "0" ]; then
  echo "!!!!!!Swagger documentation generate error, please check the source code!!!!!!"
  exit 1
fi

# build webui
#cd webui && yarn run build && cd ../
#sed -i "s/Vue App/KubeFileBrowser/g" static/index.html
# build server
name="kubefilebrowser"
# shellcheck disable=SC2006
VERSION=`git describe --tags --abbrev=0`
# shellcheck disable=SC2006
GO_VERSION=`go version|awk '{print $3" "$4}'`
# shellcheck disable=SC2006
GIT_URL=`git remote -v|grep push|awk '{print $2}'`
# shellcheck disable=SC2006
GIT_BRANCH=`git rev-parse --abbrev-ref HEAD`
# shellcheck disable=SC2006
GIT_COMMIT=`git rev-parse HEAD`
# shellcheck disable=SC2006
GIT_LATEST_TAG=`git describe --tags --abbrev=0`
# shellcheck disable=SC2006
BUILD_TIME=`date +"%Y-%m-%d %H:%M:%S %Z"`

LDFLAGS="-X 'github.com/xmapst/kubefilebrowser.Version=${VERSION}' -X 'github.com/xmapst/kubefilebrowser.GoVersion=${GO_VERSION}' -X 'github.com/xmapst/kubefilebrowser.GitUrl=${GIT_URL}' -X 'github.com/xmapst/kubefilebrowser.GitBranch=${GIT_BRANCH}' -X 'github.com/xmapst/kubefilebrowser.GitCommit=${GIT_COMMIT}' -X 'github.com/xmapst/kubefilebrowser.GitLatestTag=${GIT_LATEST_TAG}' -X 'github.com/xmapst/kubefilebrowser.BuildTime=${BUILD_TIME}'"

# linux
archList="386 amd64 arm arm64 ppc64le"
# shellcheck disable=SC2181
for i in $archList; do
  # shellcheck disable=SC2027
  BinaryName=$name"_linux-"${i}"-"${VERSION}
  # shellcheck disable=SC2090
  CGO_ENABLED=0 GOOS=linux GOARCH=$i go build -ldflags "${LDFLAGS}" -o "$BinaryName" cmd/server/main.go
  # shellcheck disable=SC2181
  if [ "$?" != "0" ]; then
    echo "!!!!!!ls compilation error, please check the source code!!!!!!"
    exit 1
  fi
  upx --lzma "$BinaryName"
done

# windows
# shellcheck disable=SC2181
archList="386 amd64"
for i in $archList; do
  # shellcheck disable=SC2027
  BinaryName=$name"_windows-"${i}"-"${VERSION}".exe"
  # shellcheck disable=SC2090
  CGO_ENABLED=0 GOOS=windows GOARCH=$i go build -ldflags "${LDFLAGS}" -o "$BinaryName" cmd/server/main.go
  # shellcheck disable=SC2181
  if [ "$?" != "0" ]; then
    echo "!!!!!!ls compilation error, please check the source code!!!!!!"
    exit 1
  fi
  upx --lzma "$BinaryName"
done

# darwin
# shellcheck disable=SC2181
archList="arm64 amd64"
for i in $archList; do
  # shellcheck disable=SC2027
  BinaryName=$name"_darwin-"${i}"-"${VERSION}
  # shellcheck disable=SC2090
  CGO_ENABLED=0 GOOS=darwin GOARCH=$i go build -ldflags "${LDFLAGS}" -o "$BinaryName" cmd/server/main.go
  # shellcheck disable=SC2181
  if [ "$?" != "0" ]; then
    echo "!!!!!!ls compilation error, please check the source code!!!!!!"
    exit 1
  fi
  upx --lzma "$BinaryName"
done