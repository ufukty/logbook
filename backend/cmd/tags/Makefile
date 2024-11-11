all: \
	database/dbtx.sqlc.go \
	client/client.go \
	client/mock.go \
	service/config.go
	$(MAKE) -C endpoints

database/dbtx.sqlc.go: database/sqlc.yml database/schema.sql
	cd database && sqlc generate .

client/client.go: $(shell find endpoints -name '*.go' ! -name '*.gh.go')
	gohandlers clients -dir endpoints -out client/client.go -pkg tags -import logbook/cmd/tags/endpoints

client/mock.go: $(shell find endpoints -name '*.go' ! -name '*.gh.go')
	gohandlers mock -dir endpoints -out client/mock.go -pkg tags -import "logbook/cmd/tags/endpoints"

service/config.go: local.yml
	gonfique generate -in local.yml -pkg service -out service/config.go