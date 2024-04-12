all: \
	database/dbtx.sqlc.go \
	service/config.go

database/dbtx.sqlc.go: database/sqlc.yml database/queries.sql database/schema.sql
	cd database && mv dbtx.sqlc.go dbtx.sqlc.go.bak && sqlc generate . ; mv dbtx.sqlc.go.bak dbtx.sqlc.go

service/config.go: testing.yml
	gonfique -in testing.yml -pkg service -out service/config.go 