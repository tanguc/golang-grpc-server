 cp ~/Projects/PersistentMarkingLB/misc/grpc/proto/upstream.proto   ../proto/ && protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative  upstream.proto
