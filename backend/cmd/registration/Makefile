all: \
	client/client.go \
	client/mock.go \
	service/config.go
	$(MAKE) -C endpoints

client/client.go: $(shell find endpoints -name '*.go' ! -name '*.gh.go')
	gohandlers clients -dir endpoints -out client/client.go -pkg registration -import logbook/cmd/registration/endpoints

client/mock.go: $(shell find endpoints -name '*.go' ! -name '*.gh.go')
	gohandlers mock -dir endpoints -out client/mock.go -pkg registration -import "logbook/cmd/registration/endpoints"

service/config.go: local.yml
	gonfique generate -in local.yml -pkg service -out service/config.go 