cd proto
protoc --go_out=. --go-grpc_out=. *
cd ..