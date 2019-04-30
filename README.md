# drone-slack-notify-log

[![Build Status](https://cloud.drone.io/api/badges/matsubara0507/drone-slack-notify-log/status.svg)](https://cloud.drone.io/matsubara0507/drone-slack-notify-log)
[![GoDoc](https://godoc.org/github.com/matsubara0507/drone-slack-notify-log?status.svg)](https://godoc.org/github.com/matsubara0507/drone-slack-notify-log)
[![Go Report Card](https://goreportcard.com/badge/github.com/matsubara0507/drone-slack-notify-log)](https://goreportcard.com/report/github.com/matsubara0507/drone-slack-notify-log)
[![](https://images.microbadger.com/badges/image/matsubara0507/slack-notify-log.svg)](https://microbadger.com/images/matsubara0507/slack-notify-log "Get your own image badge on microbadger.com")


Drone plugin for sending Drone step log to Slack as snippet.

use [files.upload slack api](https://api.slack.com/methods/files.upload).

## Build

Build the binary with the following commands:

```
go build
```

## Docker

Build the Docker image with the following commands:

```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -o release/linux/amd64/drone-slack-notify-log
docker build --rm -t matsubara0507/slack-notify-log .
```

## Usage

Execute from the working directory:

```
docker run --rm \
  -e SLACK_TOKEN=xxx \
  -e PLUGIN_CHANNEL=foo \
  -e PLUGIN_DRONE_TOKEN=yyy \
  -e PLUGIN_DRONE_HOST=https://cloud.drone.io \
  -e PLUGIN_STEP_NUMBER=1 \
  -e DRONE_REPO_OWNER=octocat \
  -e DRONE_REPO_NAME=hello-world \
  -e DRONE_COMMIT_SHA=7fd1a60b01f91b314f59955a4e4d4e80d8edf11d \
  -e DRONE_COMMIT_BRANCH=master \
  -e DRONE_COMMIT_AUTHOR=octocat \
  -e DRONE_BUILD_NUMBER=1 \
  -e DRONE_BUILD_STATUS=success \
  -e DRONE_BUILD_LINK=http://github.com/octocat/hello-world \
  -e DRONE_STAGE_NUMBER=1 \
  -e DRONE_TAG=1.0.0 \
  matsubara0507/slack-notify-log
```
