SERVICE_NAME = `echo supersvc`

build:
	@go build -ldflags "-X main.ServiceName=${SERVICE_NAME}" --race --tags=dynamic -o ./bin/api/app ./cmd/api/main.go

run: build
	@./bin/api/app