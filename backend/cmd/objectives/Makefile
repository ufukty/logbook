all: \
	database/dbtx.sqlc.go \
	client/client.go \
	service/config.go
	$(MAKE) -C endpoints

database/dbtx.sqlc.go: database/sqlc.yml database/queries.*.sql database/schema.sql
	cd database && sqlc generate .

service/config.go: local.yml
	gonfique generate -in local.yml -pkg service -out service/config.go

client/client.go: $(shell find endpoints -name '*.go' ! -name '*.gh.go')
	gohandlers clients -dir endpoints -out client/client.go -pkg objectives -import logbook/cmd/objectives/endpoints