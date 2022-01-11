#!/usr/bin/env bash

mkdir -p utils/kftoolsbinary
# linux
archList="386 amd64 arm arm64 ppc64le"
# shellcheck disable=SC2181
BuildAt=$(date)
GitHash=$(git rev-parse --short HEAD)
for i in $archList; do
  # shellcheck disable=SC2027
  BinaryName="kftools_linux_"$i
  CGO_ENABLED=0 GOOS=linux GOARCH=$i go build -a -installsuffix cgo -ldflags "-s -w -X 'main.BuildAt=$BuildAt' -X 'main.GitHash=$GitHash'" -o utils/kftoolsbinary/"$BinaryName" cmd/kftools/main.go
  # shellcheck disable=SC2181
  if [ "$?" != "0" ]; then
    echo "!!!!!!ls compilation error, please check the source code!!!!!!"
    exit 1
  fi
  upx --lzma utils/kftoolsbinary/"$BinaryName"
done

# windows
# shellcheck disable=SC2181
archList="386 amd64"
for i in $archList; do
  # shellcheck disable=SC2027
  BinaryName="kftools_windows_"$i".exe"
  CGO_ENABLED=0 GOOS=windows GOARCH=$i go build -a -installsuffix cgo -ldflags "-s -w -X 'main.BuildAt=$BuildAt' -X 'main.GitHash=$GitHash'" -o utils/kftoolsbinary/"$BinaryName" cmd/kftools/main.go
  # shellcheck disable=SC2181
  if [ "$?" != "0" ]; then
    echo "!!!!!!ls compilation error, please check the source code!!!!!!"
    exit 1
  fi
  upx --lzma utils/kftoolsbinary/"$BinaryName"
done

# darwin
# shellcheck disable=SC2181
archList="arm64 amd64"
for i in $archList; do
  # shellcheck disable=SC2027
  BinaryName="kftools_darwin_"$i
  CGO_ENABLED=0 GOOS=darwin GOARCH=$i go build -a -installsuffix cgo -ldflags "-s -w -X 'main.BuildAt=$BuildAt' -X 'main.GitHash=$GitHash'" -o utils/kftoolsbinary/"$BinaryName" cmd/kftools/main.go
  # shellcheck disable=SC2181
  if [ "$?" != "0" ]; then
    echo "!!!!!!ls compilation error, please check the source code!!!!!!"
    exit 1
  fi
  upx --lzma utils/kftoolsbinary/"$BinaryName"
done