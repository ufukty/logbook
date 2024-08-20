all: \
	queries/dbtx.sqlc.go \
	service/config.go

queries/dbtx.sqlc.go: queries/sqlc.yml queries/queries.*.sql queries/schema.sql
	cd queries && sqlc generate .

service/config.go: local.yml
	gonfique generate -in local.yml -pkg service -out service/config.go 