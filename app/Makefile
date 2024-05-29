all: app/internal/web/servdisc/models/stage/schema.go config/api/schema.go 
	$(MAKE) -C cmd/account
	$(MAKE) -C cmd/objectives

config/api/schema.go: api.yml
	gonfique \
		-pkg api \
		-mappings config/api/mappings.yml \
		-organize \
		-in api.yml \
		-out config/api/schema.go

app/internal/web/servdisc/models/stage/schema.go: ../platform/stage/artifacts/deployment/service_discovery.json
	gonfique \
		-pkg stage \
		-mappings internal/web/servdisc/models/stage/mappings.yml \
		-organize \
		-in ../platform/stage/artifacts/deployment/service_discovery.json \
		-out internal/web/servdisc/models/stage/schema.go