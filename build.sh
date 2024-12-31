#!/usr/bin/env bash

go build -ldflags='-s' -o ./dist/i3wm-obs ./
GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o ./dist/i3wm-obs_amd64 ./
GOOS=linux GOARCH=arm64 go build -ldflags='-s' -o ./dist/i3wm-obs_arm64 ./
