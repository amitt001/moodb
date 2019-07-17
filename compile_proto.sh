#!/bin/sh

protoc -I wal/walpb --go_out=plugins=grpc:wal/walpb wal/walpb/record.proto

protoc -I mdbserver/mdbserverpb --go_out=plugins=grpc:mdbserver/mdbserverpb mdbserver/mdbserverpb/command.proto