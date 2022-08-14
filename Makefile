build: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -v -tags netgo -ldflags '-w -extldflags "-static"' -o ./build/main ./cmd/main.go

clean:
	rm -r ./build