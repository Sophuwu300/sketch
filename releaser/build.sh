#!/bin/bash
command -v upx
doUpx="$?"

set -e
[[ "$TARGET" == "amd64" || "$TARGET" == "arm64" ]] || exit 1
[[ "$VERSION" == "$(date +%Y.%m.%d)" ]] || exit 1
START_DIR=$(pwd)
BUILD_DIR="$START_DIR/build/$VERSION"
SRC_DIR="$BUILD_DIR/src"
BUILD_DIR="$BUILD_DIR/$TARGET"
BIN_OUT="$BUILD_DIR/sketch"

mkdir -p $BUILD_DIR

sed -e "s/{{ Version }}/$VERSION/g" "$START_DIR/nfpm.yaml" | sed -e "s/{{ Target }}/$TARGET/g" > "$BUILD_DIR/nfpm.yaml"
cd "$BUILD_DIR"
GOOS=linux GOARCH=$TARGET go build -ldflags "-w -s" -o "$BIN_OUT" "$SRC_DIR"
if [[ $doUpx == "0" ]]; then
  upx --best --lzma "$BIN_OUT"
fi
nfpm package --packager deb --config nfpm.yaml