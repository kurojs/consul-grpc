gen: 
	protoc --proto_path=idl \
	-I$$GOPATH/src/gitlab.360live.vn/zpi \
	--go_out=plugins=grpc:./grpc-gen \
	idl/greeting.proto

build: 
	go build main.go

