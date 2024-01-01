all: \
	build/simplified.png \
	build/simplified.d2.png \
	build/versioning.v1.png \
	build/versioning.v2.png 

build:
	mkdir -p build

build/simplified.png: simplified.mmd | build
	mmdc -f -e png -s 3 -i "simplified.mmd" -o "build/simplified.png"

build/versioning.v1.png: versioning.v1.mmd | build
	mmdc -f -e png -s 3 -i "versioning.v1.mmd" -o "build/versioning.v1.png"

build/versioning.v2.png: versioning.v2.mmd | build
	mmdc -f -e png -s 3 -i "versioning.v2.mmd" -o "build/versioning.v2.png"

build/simplified.d2.png: simplified.d2 | build
	d2 --layout elk "simplified.d2" "build/simplified.d2.png"