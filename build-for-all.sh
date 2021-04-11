#!/bin/bash
GOOS=windows GOARCH=amd64 go build -o bin/techdev-amd64.exe main.go
GOOS=windows GOARCH=386 go build -o bin/techdev-386.exe main.go
GOOS=darwin GOARCH=amd64 go build -o bin/techdev-amd64-macos main.go
GOOS=linux GOARCH=amd64 go build -o bin/techdev-amd64-linux main.go