# redis-clone

A Redis server clone written in Go from scratch, no external dependencies.
Implements a raw TCP server, a wire protocol parser, and a command
dispatcher.

## Running it

```
make run
```

Starts on `:8091`. Change the port with `-port`:

```
go run ./cmd/server -port 6380
```

Then connect with nc/telnet and type commands, one per line:

```
$ nc localhost 8091
write 'exit' to exit
PING
+PONG
SET foo bar
+OK
GET foo
$3
bar
```

## Commands

`PING`, `ECHO <msg>`, `GET <key>`, `SET <key> <value>`, `EXISTS <key>`,
`DEL <key> [key ...]`, `FLUSH`. Command names are case insensitive.

## Protocol quirk

Requests are just space-separated text, one line per command. So no quoted
strings, `SET foo "two words"` just looks like too many args and errors out.

Responses are encoded in real RESP (the `+`/`-`/`:`/`$`/`*` format redis-cli
speaks). So requests use a simplified line protocol while responses use the
real wire format. A working binary RESP decoder already exists in
`internal/resp/parser.go` (`Parse`), but it isn't wired into the server's
read path yet.

## Testing

```
make test
make race
make vet
make fmt
```

Single test:

```
go test ./internal/resp -run TestParseRESP_BulkString -v
go test ./internal/command -run TestHandleSetGet_RoundTrip -v
```

## Layout

- `internal/resp` - parsing. `ParseSimple` handles the inline protocol the
  server uses for requests, `Parse`/`Encode` handle real RESP.
- `internal/command` - `Handlers.Dispatch` routes a request to a handler and
  encodes the result. Also holds the store, a map guarded by a mutex.
- `internal/server` - accepts connections, one goroutine per connection,
  each reading a request, dispatching it, and writing the response back.

## License

GPLv3, see LICENSE.
