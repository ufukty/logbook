all: \
	backend/config/deployment/config.go
	$(MAKE) -C backend
	$(MAKE) -C docs

backend/config/deployment/config.go: platform/local/deployment.yml backend/config/deployment/directives.yml
	gonfique generate \
		-pkg deployment \
		-directives backend/config/deployment/directives.yml \
		-in platform/local/deployment.yml \
		-out backend/config/deployment/config.go