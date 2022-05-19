#!/usr/bin/env bash

mkdir -p utils/kftoolsbinary
# shellcheck disable=SC2006
DistList=`go tool dist list`
# shellcheck disable=SC2181
BuildAt=$(date)
GitHash=$(git rev-parse --short HEAD)
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
  if [ "${Archs}" -ne "386" ] || [ "${Archs}" -ne "amd64" ] || [ "${Archs}" -ne "arm" ] || [ "${Archs}" -ne "arm64" ];then
    continue
  fi
  echo "Building ${Platforms} ${Archs}..."
  # shellcheck disable=SC2027
  BinaryName="kftools_${Platforms}_${Archs}"
  if [ "${Platforms}" == "windows" ];then
    BinaryName="kftools_${Platforms}_${Archs}.exe"
  fi
  CGO_ENABLED=0 GOOS=${Platforms} GOARCH=${Archs} go build -a -installsuffix cgo -ldflags "-s -w -X 'main.BuildAt=${BuildAt}' -X 'main.GitHash=${GitHash}'" -o utils/kftoolsbinary/"${BinaryName}" cmd/kftools/main.go
  # shellcheck disable=SC2181
  if [ "$?" != "0" ]; then
    echo "!!!!!!ls compilation error, please check the source code!!!!!!"
    exit 1
  fi
done