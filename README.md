# restarter
A small application that stops a Wolfenstein: Enemy Territory server after being online for X hours if the server has been empty for a while

Build for your platform:

```
https://github.com/golang/go/wiki/WindowsCrossCompiling
go build -o bin/restarter src/main.go
```

Usage:

```
./restarter -ip 123.123.123.123 -port 27965 -rconpassword foobar -interval 24h -pollrate 1m -numchecks 5
```
