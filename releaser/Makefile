# Makefile for building the deb package
NAME = sketch
# Versioning
# run `date '+%Y.%m.%d' to set the version`
VERSION = $(shell date '+%Y.%m.%d')

START_DIR = $(shell pwd)
SRC_DIR = ../src
BUILD_DIR = ./build/$(VERSION)

ifeq ($(shell command -v upx), )
UPX = echo "UPX not found, skipping compression"
else
UPX = upx --lzma --best $(NAME)
endif


FLAGS =

defualt: $(SRC_DIR)/ sketch.1
	mkdir -p $(BUILD_DIR)
	cp -r $(SRC_DIR) $(BUILD_DIR)
	sed -e 's/{{ Version }}/$(VERSION)/g' sketch.1 > ./$(BUILD_DIR)/sketch.1
	sed -e 's/{{ Version }}/$(VERSION)/g' ../src/version.go > ./$(BUILD_DIR)/src/version.go
	VERSION=$(VERSION) TARGET="amd64" ./build.sh
	VERSION=$(VERSION) TARGET="arm64" ./build.sh

clean:
	@echo "unimplemented"
