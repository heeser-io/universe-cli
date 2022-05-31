#!/bin/bash

# build for macos x64
env GOOS=darwin GOARCH=amd64 go build -o bin/macos/universe
env GOOS=darwin GOARCH=arm64 go build -o bin/macos-silicon/universe

# build for linux x64
env GOOS=linux GOARCH=amd64 go build -o bin/linux64/universe
env GOOS=linux GOARCH=386 go build -o bin/linux32/universe

# build for windows x64
env GOOS=windows GOARCH=amd64 go build -o bin/windows/universe.exe

# zip everything
# we use notarize for mac
# zip bin/macos-x64.zip bin/macos/*
# zip bin/macos-silicon.zip bin/macos-silicon/*
zip bin/linux-x64.zip bin/linux64/*
zip bin/linux-386.zip bin/linux32/*
zip bin/windows-x64.zip bin/windows/*