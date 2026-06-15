#!/usr/bin/sh

MITM_VERSION=$(git describe --tags)

go build -ldflags "-X main.version=${MITM_VERSION}" -o ./bin/mitm-collector-ora main.go

cp bin/mitm-collector-ora ../../scheduler/mitm_scheduler/bin/.