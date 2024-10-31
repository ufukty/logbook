all: \
	config/api/schema.go
	$(MAKE) -C cmd
	$(MAKE) -C models

config/api/schema.go: api.yml config/api/directives.yml
	gonfique generate \
		-pkg api \
		-directives config/api/directives.yml \
		-in api.yml \
		-out config/api/schema.go