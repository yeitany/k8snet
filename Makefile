#test:
#	go test ./...
lint:
	go fmt ./...
	go vet ./...

build_path=./.bin/
build-binaries:
	GOOS=linux GOARC=amd64 go build -o $(build_path)/main ./cmd/k8snet/main.go

build-charts:
	helm package charts/k8snet