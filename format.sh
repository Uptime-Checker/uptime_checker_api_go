#! /bin/sh

golines . -m 120 -w
golangci-lint run
