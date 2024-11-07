all: bindings.gh.go list.gh.go

list.gh.go: $(shell find . -name '*.go' ! -name '*.gh.go')
	gohandlers list -dir .

bindings.gh.go: $(shell find . -name '*.go' ! -name '*.gh.go')
	gohandlers bindings -dir . -recv Private