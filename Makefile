.PHONY: protos

protos:
	 protoc proto/ipcamera.proto --go_out=plugins=grpc:proto/ipcamera