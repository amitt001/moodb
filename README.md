# MooDB(mdb)

It's a WIP(Work in progress)

A simple persistent key-value database implemeneted in Go.

Blog: https://kakku.org/writing-a-simple-database/


## Run

1. Edit config/server.yaml and config/client.yaml files to put right value for WAL datadir

2. Server: `go run cmd/server/main.go -logtostderr=true`

3. Client: `go run mdbcli/*.go`

## Commands

In the client shell

```
⇒  go run mdbcli/*.go
```

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

## Log compaction

- Sync policy
- Compact by file size or by percentage increase
- Truncate file at startup
- Generate a snapshot file from loaded data
