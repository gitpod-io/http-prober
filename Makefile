build: http-prober

# Build api binary
http-prober: fmt vet
	CGO_ENABLED=0 go build -v -ldflags '-w -extldflags '-static'' -o http-prober

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...