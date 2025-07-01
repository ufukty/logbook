all: client/client.go
	$(MAKE) -C endpoints

client/client.go: $(shell find endpoints -name '*.go' ! -name '*.gh.go')
	gohandlers client -dir endpoints -out client/client.go -pkg pdp -import logbook/cmd/pdp/endpoints
