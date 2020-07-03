.PHONY: compile

compile:
	protoc -Iproto --go_out=plugins=grpc:api proto/*.proto
