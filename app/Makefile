all: \
	internal/web/discovery/models/stage/stage.go \
	internal/web/discovery/models/local/local.go \
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

internal/web/discovery/models/stage/stage.go: internal/web/discovery/models/stage/service_discovery.json internal/web/discovery/models/stage/directives.yml
	gonfique generate \
		-pkg stage \
		-directives internal/web/discovery/models/stage/directives.yml \
		-type-name "Stage" \
		-in internal/web/discovery/models/stage/service_discovery.json \
		-out internal/web/discovery/models/stage/stage.go

internal/web/discovery/models/local/local.go: internal/web/discovery/models/local/service_discovery.yml internal/web/discovery/models/local/directives.yml
	gonfique generate \
		-pkg local \
		-directives internal/web/discovery/models/local/directives.yml \
		-type-name "Local" \
		-in internal/web/discovery/models/local/service_discovery.yml \
		-out internal/web/discovery/models/local/local.go

config/deployment/config.go: ../platform/testing.yml config/deployment/directives.yml
	gonfique generate \
		-pkg deployment \
		-directives config/deployment/directives.yml \
		-in ../platform/testing.yml \
		-out config/deployment/config.go