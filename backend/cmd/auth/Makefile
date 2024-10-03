all: web/build

web/build: web/src/**
	rm -rfv web/build
	kask build -in web/src -out web/build -domain "/"

watch:
	fswatch web/src | while read l; do $(MAKE) -B web/build; done