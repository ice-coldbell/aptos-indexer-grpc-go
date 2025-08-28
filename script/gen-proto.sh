#!/bin/bash

protoc \
--experimental_allow_proto3_optional \
--proto_path=proto \
--go_out=. \
--go_opt=paths=source_relative \
--go-grpc_out=. \
--go-grpc_opt=paths=source_relative \
    proto/aptos/indexer/v1/grpc.proto \
    proto/aptos/indexer/v1/raw_data.proto \
    proto/aptos/indexer/v1/filter.proto \
    proto/aptos/transaction/v1/transaction.proto \
    proto/aptos/util/timestamp/timestamp.proto \
    proto/aptos/remote_executor/v1/network_msg.proto \
    proto/aptos/internal/fullnode/v1/fullnode_data.proto