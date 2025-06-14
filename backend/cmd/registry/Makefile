all: client/client.go client/mock.go
	$(MAKE) -C endpoints

client/client.go: $(shell find endpoints -name '*.go' ! -name '*.gh.go')
	gohandlers client -dir endpoints -out client/client.go -pkg registry -import logbook/cmd/registry/endpoints

client/mock.go: $(shell find endpoints -name '*.go' ! -name '*.gh.go')
	gohandlers mock -dir endpoints -out client/mock.go -pkg registry -import "logbook/cmd/registry/endpoints"