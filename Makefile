SOURCES=./

dep:
	dep ensure

http_example:
	go run ${SOURCES}/examples/http/main.go

.DEFAULT_GOAL := test
test: http_example

