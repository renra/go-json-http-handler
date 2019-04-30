SOURCES=./

dep:
	dep ensure

http_example:
	go run ${SOURCES}/examples/http/main.go

.PHONY: test
.DEFAULT_GOAL := test
test:
	go test ./test... -count 1 -v
