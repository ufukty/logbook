FILES := $(shell find . -type f | grep -v ".touch.make")

build: .touch.make

.touch.make: $(FILES)
	bash build.sh
	@touch $@
