all: \
	config/api/schema.go \
	config/deployment/config.go
	$(MAKE) -C cmd/account
	$(MAKE) -C cmd/objectives

config/api/schema.go: api.yml config/api/directives.yml
	gonfique generate \
		-pkg api \
		-directives config/api/directives.yml \
		-in api.yml \
		-out config/api/schema.go

config/deployment/config.go: ../platform/local/deployment.yml config/deployment/directives.yml
	gonfique generate \
		-pkg deployment \
		-directives config/deployment/directives.yml \
		-in ../platform/local/deployment.yml \
		-out config/deployment/config.go