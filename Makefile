all: build

build: build_proto

build_proto: 
	protoc --go_out=. ./message/message.proto