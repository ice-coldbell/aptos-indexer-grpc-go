# aptos-indexer-grpc-go


# Install dependency

## [protoc](https://protobuf.dev/installation/)
- Linux, using apt or apt-get, for example:
```bash
apt install -y protobuf-compiler
protoc --version  # Ensure compiler version is 3+
```
- MacOS, using Homebrew:
```zsh
brew install protobuf
protoc --version  # Ensure compiler version is 3+
```

## [protoc-gen-go, protoc-gen-go-grpc](https://grpc.io/docs/languages/go/quickstart/)
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
```

# Build
```bash
./script/gen-proto.sh
```