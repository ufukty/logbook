all: handlers.go

handlers.go: $(shell find . -name '*.go' ! -name 'handlers.go')
	gohandlers -dir .
