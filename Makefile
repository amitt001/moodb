build:
	mkdir -p data
	bash ./compile_proto.sh
	go build ./...
