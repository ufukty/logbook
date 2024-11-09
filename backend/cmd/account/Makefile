all: \
	database/dbtx.sqlc.go \
	service/config.go
	$(MAKE) -C endpoints

database/dbtx.sqlc.go: database/sqlc.yml database/queries.sql database/schema.sql
	cd database && sqlc generate .

service/config.go: local.yml
	gonfique generate -in local.yml -pkg service -out service/config.go 