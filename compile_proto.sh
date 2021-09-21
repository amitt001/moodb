#!/bin/sh

protoc --proto_path=mdbserver/mdbserverpb --go_out=mdbserver/mdbserverpb --go_opt=paths=source_relative --go-grpc_out=mdbserver/mdbserverpb --go-grpc_opt=paths=source_relative mdbserver/mdbserverpb/*.proto