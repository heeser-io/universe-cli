#!/bin/bash
git describe --tags

VERSION=$(git describe --tags)
docker buildx build --platform linux/amd64 -t universe-cli .
docker tag universe-cli:latest ghcr.io/heeser-io/universe-cli:$VERSION
docker push ghcr.io/heeser-io/universe-cli:$VERSION