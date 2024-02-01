all: internal/web/api/schema.go

internal/web/api/schema.go: api.yml
	gonfique -in api.yml -out internal/web/api/schema.go -pkg api

