# EchoServer                                                                    
simple echo server for test

## Build
1. clone the source under $GOPATH/src
2. Enter $GOPATH/src/EchoServer directory, run 'go build EchoServer.go'

## Execute

### Run directly
```
./EchoServer
```
it will "Listen on Port=8888, Timeout=10m0s, EchoMode=Single" by default

```
Usage of ./EchoServer:
  -b  Use broadcast mode
  -p  Set Listen Port, default 8888 
  -t  Set Timeout, default 600 (seconds), that is, 10minutes
```

## Echo Mode

### Single Mode
When echo server receives data from a client, it sends the data back to the client directly and other clients received no data. 

### Broadcast Mode
When echo server receives data from a client, it broadcasts the data to all the clients connected. 
