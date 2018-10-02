#!/usr/bin/env bash

protoc -I${GOOGLEAPIS_DIR}  -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
-I=./ --go_out=plugins=grpc:. bclogpb.proto
