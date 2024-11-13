all: \
	client/client.go \
	client/mock.go \
	service/config.go
	$(MAKE) -C endpoints

client/client.go: $(shell find endpoints -name '*.go' ! -name '*.gh.go')
	gohandlers clients -dir endpoints -out client/client.go -pkg account -import logbook/cmd/account/endpoints

client/mock.go: $(shell find endpoints -name '*.go' ! -name '*.gh.go')
	gohandlers mock -dir endpoints -out client/mock.go -pkg account -import "logbook/cmd/account/endpoints"

service/config.go: local.yml
	gonfique generate -in local.yml -pkg service -out service/config.go 