#!/bin/bash

GOARCH=amd64

GOOS=darwin
go build -o terraform-provider-n0stack-yamlcompiler-$GOOS-$GOARCH

GOOS=linux
go build -o terraform-provider-n0stack-yamlcompiler-$GOOS-$GOARCH

GOOS=windows
go build -o terraform-provider-n0stack-yamlcompiler-$GOOS-$GOARCH
