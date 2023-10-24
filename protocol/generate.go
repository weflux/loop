package protocol

//go:generate protoc -I ./third_party -I . --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ./shared/*.proto ./loopify/v1/*.proto ./proxy/*.proto ./broker/*.proto
