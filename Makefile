cur_dir=$(shell pwd)
bin_dir=$(cur_dir)/bin/
BUILD_GCFLAG="-N -l"
$(shell go env -w GO111MODULE=on)
$(shell go env -w GOFLAGS=-mod=vendor)

BUILD_CMD=go build
BUILD_FLAGS=-gcflags $(BUILD_GCFLAG)  -trimpath -a -o $(bin_dir)

.PHONY: vendor
vendor:
	rm -rf go.sum vendor vendor.tar.gz \
	export GOPROXY=https://goproxy.io,direct \
	&& go mod tidy \
	&& go mod vendor \
	&& tar -czf vendor.tar.gz vendor


deps:
	rm -rf vendor && tar -xzf vendor.tar.gz


bin: deps
	$(BUILD_CMD) $(BUILD_FLAGS)/htools ./cmd
