all: web/build

web/build: web/kask/** web/ts/src/**
	rm -rfv web/build
	kask build -in web/kask -out web/build -domain "/"
	cd web/ts && tsc

watch:
	fswatch web/kask | while read l; do $(MAKE) -B web/build; done
