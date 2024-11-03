all: build.gh.go list.gh.go

list.gh.go: $(shell find . -name '*.go' ! -name '*.gh.go')
	gohandlers list -dir .

build.gh.go: $(shell find . -name '*.go' ! -name '*.gh.go')
	gohandlers build -dir . -recv Private