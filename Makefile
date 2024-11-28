generate_proto:
	for proto_file in $(wildcard api/*.proto); do \
	  protoc -I. -Iapi/googleapis --go_out=pkg --go-grpc_out=pkg --grpc-gateway_out=pkg $$proto_file --experimental_allow_proto3_optional; \
	done

init_db:
	sh scripts/create-tables.sh

drop_db:
	sh scripts/drop-tables.sh
