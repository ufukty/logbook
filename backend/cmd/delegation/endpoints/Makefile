all: bindings.gh.go validate.gh.go list.gh.go

list.gh.go: $(shell find . -name '*.go' ! -name '*.gh.go')
	gohandlers list -dir . -hi-type reception.HandlerInfo -hi-import logbook/internal/web/reception

validate.gh.go: $(shell find . -name '*.go' ! -name '*.gh.go')
	gohandlers validate -dir .

bindings.gh.go: $(shell find . -name '*.go' ! -name '*.gh.go')
	gohandlers bindings -dir .