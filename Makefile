all: build

build: build_proto

build_proto: 
	protoc --go_out=. ./message/message.proto

run_tests:
	go test ./base58 ./wallet ./blockchain -tags=tests
