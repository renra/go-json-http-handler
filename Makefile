SOURCES=./

dep:
	dep ensure

example:
	go run ${SOURCES}/examples/main.go

.DEFAULT_GOAL := test
test: example

