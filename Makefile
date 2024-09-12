all: app/config/deployment/config.go
	$(MAKE) -C app
	$(MAKE) -C docs
