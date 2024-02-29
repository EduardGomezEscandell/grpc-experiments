#!/bin/bash
set -eu

cd hello/

PATH=$PATH:$(go env GOPATH)/bin protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative hello.proto

python3 -m grpc_tools.protoc -I. --python_out=. --pyi_out=. --grpc_python_out=. --experimental_allow_proto3_optional hello.proto
mv *.pyi ..
mv *.py ..


cd ..