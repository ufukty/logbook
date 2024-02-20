all: internal/web/api/schema.go

internal/web/api/schema.go: api.yml
	gonfique \
		-pkg api \
		-mappings internal/web/api/mappings.yml \
		-organize \
		-in api.yml \
		-out internal/web/api/schema.go

