# csv-writer

An exporter that writes records to **any** `io.Writer` — a file, a buffer, a
gzip stream, an HTTP response, a test double.

## API

### `WriteRecords(w io.Writer, recs []Record) error`

Writes a CSV document with the header `name,email`, then one line per record.

1. The output is RFC 4180 CSV: a field containing a comma, a double quote or a
   newline is wrapped in double quotes, and an embedded `"` is doubled.
2. The header is always written, even when `recs` is empty.
3. The returned error is the writer's error. A failed write must never be
   reported as success, and a short write is an error.
4. `WriteRecords` does not close, flush or wrap `w`: if `w` is a `*bufio.Writer`
   or a `gzip.Writer`, the caller flushes it.
5. `Record` is:

```go
type Record struct {
	Name  string
	Email string
}
```

## Tests

```
go test ./...
```
