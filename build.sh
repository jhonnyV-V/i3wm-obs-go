#!/usr/bin/env bash

go build -ldflags='-s' -o ./dist/phoemux ./
GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o ./dist/phoemux_amd64 ./
GOOS=linux GOARCH=arm64 go build -ldflags='-s' -o ./dist/phoemux_arm64 ./
