#!/bin/bash

docker buildx build --platform linux/amd64 -t universe-cli .
docker tag universe-cli:latest yhc44/universe-cli:staging
docker push yhc44/universe-cli:staging