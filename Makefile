all: app/config/deployment/config.go
	$(MAKE) -C app

app/config/deployment/config.go: platform/testing.yml
	gonfique -in platform/testing.yml -out app/config/deployment/config.go -pkg deployment