## Install protoc
1. Download the latest **Protocol Compiler (protoc)** from [release](https://github.com/protocolbuffers/protobuf/releases/tag/v25.3)
2. Unzip the file and move the protoc binary to /usr/local/bin
```bash
$ sudo mv protoc /usr/local/bin
```
3. Install the Go Protocol Buffers plugin
```bash
$ go get -u github.com/golang/protobuf/protoc-gen-go
```
4. Add the Go bin directory to your system's PATH by running `nano ~/.zshrc` to open the file and adding the following line at the end of the file:
```bash
$ export PATH=$PATH:(go env GOPATH)/bin
```
5. Reopen terminal and verify that the installation was successful
```bash
$ protoc-gen-go --version
```
## Create a .proto file
Specify the Go package within the .proto file.
```proto
syntax = "proto3";

package geecachepb;

// Specify the Go package where the code will be generated.
option go_package = "/";

message Request {
  string group = 1;
  string key = 2;
}

message Response {
  bytes value = 1;
}

service GroupCache {
  rpc Get(Request) returns (Response);
}
```
## Generate Go code
```bash
$ protoc --go_out=. *.proto
```

