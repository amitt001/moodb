# MooDB(mdb)

A simple persistent key-value database implemented in Go.

Offers a simple Write-Ahead-Logging (WAL) implementation that allows data restore across session.

Blog: https://kakku.org/writing-a-simple-database/

## Run
### Prerequisite

- Make sure you have protocolbuffer installed
  - https://grpc.io/docs/protoc-installation/
  - Verify `protoc --help`
- Make sure Go 1.17.x is installed
  - https://golang.org/dl/
  - Verify `go version`
- Make sure you have Go installed and $GOPATH set
  - Check `echo $GOPATH` returns a valid path
- Clone repository `git clone https://github.com/amitt001/moodb.git`

- `cd moodb`

1. Run build: `make build`
2. Run server: `go run cmd/server/main.go -logtostderr=true`
3. Open a new terminal window and run client: `go run mdbcli/*.go`

## Commands

In the client shell

```
MooDB version 0.0.1
o> set name Amit
Inserted 1
o> get name
Amit
o> del name
Deleted 1
o> get name
```

## WAL

To just install write-ahead log Wal module can be used separately.

`go get "github.com/amitt001/moodb/wal"`

