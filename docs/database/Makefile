all: \
	build/simplified.png \
	build/versioning.v1.png \
	build/versioning.v2.png 

build:
	mkdir -p build

build/simplified.png: simplified.mmd | build
	mmdc -f -e png -w 5000 -H 5000 -s 4 -b "#222" -t dark -i "simplified.mmd" -o "build/simplified.png"

build/versioning.v1.png: versioning.v1.mmd | build
	mmdc -f -e png -w 5000 -H 5000 -s 4 -b "#222" -t dark -i "versioning.v1.mmd" -o "build/versioning.v1.png"

build/versioning.v2.png: versioning.v2.mmd | build
	mmdc -f -e png -w 5000 -H 5000 -s 4 -b "#222" -t dark -i "versioning.v2.mmd" -o "build/versioning.v2.png"