all: \
	client/client.go \
	service/config.go
	$(MAKE) -C endpoints

client/client.go: $(shell find endpoints -name '*.go' ! -name '*.gh.go')
	gohandlers client -dir endpoints -out client/client.go -pkg registration -import logbook/cmd/registration/endpoints

service/config.go: local.yml
	gonfique generate -in local.yml -pkg service -out service/config.go 