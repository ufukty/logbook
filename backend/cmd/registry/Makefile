all: \
	client/client.go
	$(MAKE) -C endpoints

client/client.go: $(shell find endpoints -name '*.go' ! -name '*.gh.go')
	gohandlers clients -dir endpoints -out client/client.go -pkg registry -import logbook/cmd/registry/endpoints