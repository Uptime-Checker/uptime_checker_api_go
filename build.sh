#! /bin/sh

if ! [ -x "$(command -v go)" ]; then
    echo "go is not installed"
    exit 1
fi
if ! [ -x "$(command -v git)" ]; then
    echo "git is not installed"
    exit 1
fi
if [ -z "${GOPATH}" ]; then
    echo "set GOPATH"
    exit 1
fi

go mod verify
go mod vendor
go mod tidy

./format.sh
go build -o bin
