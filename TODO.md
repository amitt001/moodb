# TODO

[x] In-memory key-value store

[x] Locks

[x] DB client and server

[x] Config

[x] WAL

[x] WAL recovery


[] WAL file lock

[] In-memory db snapshot

[] Authentication

[] Logging

[] Error check. Is there a better approach then using check() method?

[] Client sends an unique id. Keep that in server, validate?

## Improvements

[] Create a .tmp directory and use .tmp files instead of writing to the file directly

[] Check how to run cleanup tasks like closing wal when go program ends