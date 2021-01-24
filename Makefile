# -------------------------------------------------
#  Pre-Release
# -------------------------------------------------
prep:
	go mod tidy &&\
	go vet ./... &&\
	go fmt ./... &&\
	go test `go list ./...`

build: prep
	GOOS=linux go build src/main.go && \
	zip main.zip main

clean:
	rm -f main main.zip

