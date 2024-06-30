all: internal/web/discovery/stage.go config/api/schema.go config/deployment/config.go
	$(MAKE) -C cmd/account
	$(MAKE) -C cmd/objectives

config/api/schema.go: api.yml config/api/directives.yml
	gonfique generate \
		-pkg api \
		-directives config/api/directives.yml \
		-in api.yml \
		-out config/api/schema.go

internal/web/discovery/stage.go: internal/web/discovery/testdata/service_discovery_stage.json internal/web/discovery/stage-mapping.yml
	gonfique \
		-pkg discovery \
		-mappings internal/web/discovery/stage-mapping.yml \
		-type-name "Stage" \
		-in internal/web/discovery/testdata/service_discovery_stage.json \
		-out internal/web/discovery/stage.go

config/deployment/config.go: ../platform/testing.yml config/deployment/directives.yml
	gonfique generate \
		-pkg deployment \
		-directives config/deployment/directives.yml \
		-in ../platform/testing.yml \
		-out config/deployment/config.go