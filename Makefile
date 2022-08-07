#test:
#	go test ./...
build_path="./.bin/"
build:
	GOOS=linux GOARC=amd64 go build -o $(build_path)/conntrack ./cmd/conntrack.go