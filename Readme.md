# Trying out gRPC with Go

Testing basic functionalities with gRPC


generate protobuf:
```ssh
$ protoc --go_out=. --go-grpc_out=. proto/chat.proto
```

if you can't resolve the PATH, you can:

```ssh
$ export GO_PATH=~/go

$ export PATH=$PATH:/$GO_PATH/bin
```