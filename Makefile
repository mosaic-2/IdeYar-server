generate_proto:
	for proto_file in $(wildcard api/*.proto); do \
	  protoc -I. -Iapi/googleapis --go_out=pkg --go-grpc_out=pkg --grpc-gateway_out=pkg $$proto_file; \
	done

init_db:
	sh scripts/create-tables.sh

drop_db:
	sh scripts/drop-tables.sh
