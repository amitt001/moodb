# WAL: Write Ahead Log

How data travels?

1. **Application Buffer** or **Library Buffer**: at this point data lives in application address space.
2. **Kernel Buffer**: Kernel page cache where data can live for a non fixed time.
3. **Storage Device Cache**:
4. **Stable Storage**:

How data travels in MDB?

When data is written to the network socket it is read by the server in buffer i.e. application memory. The server already has a file stream opened(WAL file).

## Installation

`go get -u github.com/tidwall/wal`

## Example

```
// Start a new WAL
walObj, err = wal.Open(/path/to/wal/dir)

// Open the existing latest WAL
walObj, err = wal.Open(/path/to/wal/dir)

// Write to WAL
err := walObj.Write([]byte("some data"))

// Read opened WAL file:
rChan := db.rWalObj.Read()
for record := range rChan {
    // do something with the data here
}
```

- Close the opened WAL:

db.rWalObj.Close()

