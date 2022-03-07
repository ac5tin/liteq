test: 
	go test ./... -v -cover

build:
	go build -o bin/app

run:
	go run main.go

run-example:
	go run example/main.go

clean:
	rm -rf ./bin

protoc:
	protoc --go_out=. --go-grpc_out=. ./queue/proto/queue.proto

build-podman:
	podman build --platform linux/amd64,linux/arm64 --format docker --manifest liteq .

build-docker:
	docker buildx build --platform linux/amd64,linux/arm64 -t docker.io/ac5tin/liteq . --push

push-podman:
	podman login docker.io
	podman manifest push --format v2s2 liteq docker://docker.io/ac5tin/liteq