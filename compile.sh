#!/bin/bash
mkdir -p build/macOS_intel/
mkdir -p build/win/
mkdir -p build/linux_intel/
mkdir -p build/linux_arm32/
mkdir -p build/linux_arm64/
mkdir -p build/freebsd/
export GOOS=darwin
export GOARCH=amd64
go build -o build/macOS_intel/
export GOOS=darwin
export GOARCH=arm64
go build -o build/macOS_m1/
export GOOS=windows
export GOARCH=amd64
go build -o build/win/
export GOOS=linux
export GOARCH=amd64
go build -o build/linux_intel/
export GOOS=linux
export GOARCH=arm64
go build -o build/linux_arm64/
export GOOS=linux
export GOARCH=arm
export GOARM=7
go build -o build/linux_arm32/
export GOOS=freebsd
export GOARCH=amd64
go build -o build/freebsd/
chmod -R +x build/
