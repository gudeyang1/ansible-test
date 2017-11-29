#!/bin/bash
rm -rf suitectl suitectl.exe

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o suitectl.exe main.go
go build -o suitectl main.go
