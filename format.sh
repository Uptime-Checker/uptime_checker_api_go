#! /bin/sh


golangci-lint run
golines . -m 120 -w
