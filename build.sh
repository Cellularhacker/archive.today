#!/bin/bash

CURRENT_DIR='./bin/current'
OS = (windows linux darwin freebsd)

mkdir -p "${CURRENT_DIR}"
mkdir -p "${CURRENT_DIR}/"
GOOS=linux GOARCH=amd64 go build -v archive.today &&