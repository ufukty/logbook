all: client/client.go client/mock.go
	$(MAKE) -C endpoints

client/client.go: $(shell find endpoints -name '*.go' ! -name '*.gh.go')
	gohandlers client -dir endpoints -out client/client.go -pkg pdp -import logbook/cmd/pdp/endpoints

client/mock.go: $(shell find endpoints -name '*.go' ! -name '*.gh.go')
	gohandlers mock -dir endpoints -out client/mock.go -pkg pdp -import "logbook/cmd/pdp/endpoints"