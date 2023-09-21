package protocol

//go:generate protoc -I ./third_party -I . --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ./shared/*.proto ./api/v1/*.proto ./message/v1/*.proto ./proxy/*.proto ./broker/*.proto
