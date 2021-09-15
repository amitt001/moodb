# MooDB(mdb)

It's a WIP(Work in progress)

A simple persistent key-value database implemented in Go.

Blog: https://kakku.org/writing-a-simple-database/

## Run

- Make sure you have Go installed and $GOPATH set
  - Check `echo $GOPATH` returns a valid path

- Run `mkdir -p $GOPATH/src/github.com/amitt001`

- Clone repository `git clone https://github.com/amitt001/moodb.git`

- `cd $GOPATH/src/github.com/amitt001/moodb`

- Generate protobuf client and server
  - `protoc --proto_path=mdbserver/mdbserverpb --go_out=mdbserver/mdbserverpb --go_opt=paths=source_relative --go-grpc_out=mdbserver/mdbserverpb --go-grpc_opt=paths=source_relative mdbserver/mdbserverpb/*.proto`

- Create data directory to save write-ahead-log(WAL) files `mkdir data`

1. Run server: `go run cmd/server/main.go -logtostderr=true`

2. Open a new terminal window and run client: `go run mdbcli/*.go`

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

Wal module can be used separately.

`go get "github.com/amitt001/moodb/wal"`

**Usage**:

```
// Existing WAL
walObj, err := wal.Open("/path/to/dir")
for i := range walObj.Read() {
			fmt.Print(i)
}

// New WAL
walObj, err = wal.New(dirPath)
```

## Log compaction

- Sync policy
- Compact by file size or by percentage increase
- Truncate file at startup
- Generate a snapshot file from loaded data

## Debug

```
>dlv debug mdbcli/*.go

set breakpoint <file>:<line number>
> break main.go:8
list
> l
start step
> s
```
