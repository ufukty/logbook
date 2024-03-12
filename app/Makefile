all: config/api/schema.go

config/api/schema.go: api.yml
	gonfique \
		-pkg api \
		-mappings config/api/mappings.yml \
		-organize \
		-in api.yml \
		-out config/api/schema.go

