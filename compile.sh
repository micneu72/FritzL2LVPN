mkdir -p build/macOS_intel/
mkdir -p build/win/
mkdir -p build/linux_intel/
mkdir -p build/linux_arm/
mkdir -p build/freebsd/
export GOOS=darwin
export GOARCH=amd64
go build -o build/macOS_intel/
export GOOS=windows
export GOARCH=amd64
go build -o build/win/
export GOOS=linux
export GPARCH=amd64
go build -o build/linux_intel/
export GOOS=linux
export GPARCH=arm
export GOARM=5
go build -o build/linux_arm/
export GOOS=freebsd
export GOARCH=amd64
go build -o build/freebsd/