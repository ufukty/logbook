all: \
	app/config/deployment/config.go
	$(MAKE) -C app
	$(MAKE) -C docs

app/config/deployment/config.go: platform/local/deployment.yml app/config/deployment/directives.yml
	gonfique generate \
		-pkg deployment \
		-directives app/config/deployment/directives.yml \
		-in platform/local/deployment.yml \
		-out app/config/deployment/config.go