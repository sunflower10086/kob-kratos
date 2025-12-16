API_PROTO_FILES=$(shell find ../../api/protos/matching -name *.proto)

.PHONY: server
server:
	kratos proto server $(API_PROTO_FILES) -t internal/service

	
