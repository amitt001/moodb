#!/bin/sh

protoc -I mdbserver/mdbserverpb --go_out=plugins=grpc:mdbserver/mdbserverpb mdbserver/mdbserverpb/command.proto