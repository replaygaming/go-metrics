default:
	mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -v -o bin/metrics_darwin_amd64
	GOOS=linux GOARCH=amd64 go build -v -o bin/metrics_linux_amd64
