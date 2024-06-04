all: internal/web/discovery/stage.go config/api/schema.go 
	$(MAKE) -C cmd/account
	$(MAKE) -C cmd/objectives

config/api/schema.go: api.yml config/api/mappings.yml
	gonfique \
		-pkg api \
		-mappings config/api/mappings.yml \
		-in api.yml \
		-out config/api/schema.go

internal/web/discovery/stage.go: internal/web/discovery/testdata/service_discovery_stage.json internal/web/discovery/stage-mapping.yml
	gonfique \
		-pkg discovery \
		-mappings internal/web/discovery/stage-mapping.yml \
		-type-name "Stage" \
		-in internal/web/discovery/testdata/service_discovery_stage.json \
		-out internal/web/discovery/stage.go
