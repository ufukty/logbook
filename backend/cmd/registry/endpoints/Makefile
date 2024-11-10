all: build.gh.go list.gh.go

list.gh.go: $(shell find . -name '*.go' ! -name '*.gh.go')
	gohandlers list -dir . -hi-type reception.HandlerInfo -hi-import logbook/internal/web/reception

build.gh.go: $(shell find . -name '*.go' ! -name '*.gh.go')
	gohandlers bindings -dir .