all: internal/web/discovery/stage.go config/api/schema.go 
	$(MAKE) -C cmd/account
	$(MAKE) -C cmd/objectives

config/api/schema.go: api.yml
	gonfique \
		-pkg api \
		-mappings config/api/mappings.yml \
		-organize \
		-in api.yml \
		-out config/api/schema.go

internal/web/discovery/stage.go: ../platform/stage/artifacts/deployment/service_discovery.json internal/web/discovery/stage-mapping.yml
	gonfique \
		-pkg discovery \
		-mappings internal/web/discovery/stage-mapping.yml \
		-type-name "Stage" \
		-in ../platform/stage/artifacts/deployment/service_discovery.json \
		-out internal/web/discovery/stage.go