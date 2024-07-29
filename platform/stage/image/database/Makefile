BASE_IMAGE_ID = $(shell doctl compute image list | grep packer-internal | tail -n 1 | cut -d ' ' -f 1)
VPC_UUID = $(shell doctl vpcs list | grep dth-fra1 | cut -d ' ' -f 1)
FILES := $(shell find . -type f | grep -v ".touch.make")

build: .touch.make

.touch.make: $(FILES)
	BASE_IMAGE_ID=$(BASE_IMAGE_ID) VPC_UUID=$(VPC_UUID) packer build -force -on-error=abort database.pkr.hcl
	@touch $@
