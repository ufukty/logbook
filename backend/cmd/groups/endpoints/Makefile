all: handlers.go

handlers.go: $(shell find . -name '*.go' ! -name 'handlers.go')
	gohandlers -dir . -type "reception.HandlerInfo" -import "logbook/internal/web/reception"