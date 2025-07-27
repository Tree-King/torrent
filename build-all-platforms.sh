#!/usr/bin/env bash
set -euo pipefail
PLATFORMS=(
  linux/amd64
  linux/arm64
  windows/amd64
  windows/arm64
  darwin/amd64
  darwin/arm64
)
OUTPUT_DIR=build
mkdir -p "$OUTPUT_DIR"
for cmd in cmd/*; do
  NAME=$(basename "$cmd")
  for PLATFORM in "${PLATFORMS[@]}"; do
    IFS=/ read GOOS GOARCH <<<"$PLATFORM"
    EXT=""
    if [[ "$GOOS" == "windows" ]]; then
      EXT=.exe
    fi
    OUT_DIR="$OUTPUT_DIR/$NAME/${GOOS}_${GOARCH}"
    mkdir -p "$OUT_DIR"
    echo "Building $NAME for $GOOS/$GOARCH"
    env CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" \
      go build -o "$OUT_DIR/$NAME$EXT" "./cmd/$NAME"
  done
done

