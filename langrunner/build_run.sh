#!/bin/bash
set -oue
docker build . -f Dockerfile.cli -t runnercli
docker run runnercli:latest
