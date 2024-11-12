all: services.validate.go
	$(MAKE) -C basics
	$(MAKE) -C columns
	$(MAKE) -C incoming

services.validate.go: services.go
	govalid -in services.go -type Service