all: database/dbtx.sqlc.go

database/dbtx.sqlc.go: database/sqlc.yml database/queries.sql database/schema.sql
	cd database && mv dbtx.sqlc.go dbtx.sqlc.go.bak && sqlc generate . ; mv dbtx.sqlc.go.bak dbtx.sqlc.go
