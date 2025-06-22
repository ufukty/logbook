all: services.validate.go
	$(MAKE) -C transports

services.validate.go: services.go
	govalid -in services.go -type Service