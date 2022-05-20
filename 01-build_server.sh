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

LDFLAGS="-X 'github.com/xmapst/kubefilebrowser.GoVersion=${GO_VERSION}' -X 'github.com/xmapst/kubefilebrowser.GitUrl=${GIT_URL}' -X 'github.com/xmapst/kubefilebrowser.GitBranch=${GIT_BRANCH}' -X 'github.com/xmapst/kubefilebrowser.GitCommit=${GIT_COMMIT}' -X 'github.com/xmapst/kubefilebrowser.BuildTime=${BUILD_TIME}'"
# shellcheck disable=SC2006
DistList=`go tool dist list`
for i in ${DistList}; do
  # shellcheck disable=SC2006
  Platforms=`echo "${i}"|awk -F/ '{print $1}'`
  if [ "${Platforms}" == "android" ] || [ "${Platforms}" == "ios" ] || [ "${Platforms}" == "js" ];then
    continue
  fi
  # shellcheck disable=SC2006
  # shellcheck disable=SC2034
  Archs=`echo "${i}"|awk -F/ '{print $2}'`
  # shellcheck disable=SC2170
  # shellcheck disable=SC2252
  if [ "${Archs}" != "386" ] && [ "${Archs}" != "amd64" ] && [ "${Archs}" != "arm" ] && [ "${Archs}" != "arm64" ];then
    continue
  fi
  echo "Building ${Platforms} ${Archs}..."
  # shellcheck disable=SC2027
  BinaryName="kftools_${Platforms}_${Archs}"
  if [ "${Platforms}" == "windows" ];then
    BinaryName="kftools_${Platforms}_${Archs}.exe"
  fi
  CGO_ENABLED=0 GOOS=${Platforms} GOARCH=${Archs} go build -a -ldflags "${LDFLAGS}" -o "${BinaryName}" cmd/server/main.go
  # shellcheck disable=SC2181
  if [ "$?" != "0" ]; then
    echo "!!!!!!ls compilation error, please check the source code!!!!!!"
    exit 1
  fi
done
