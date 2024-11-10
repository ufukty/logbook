all: services.validate.go
	$(MAKE) -C columns

services.validate.go: services.go
	govalid -in services.go -type Service