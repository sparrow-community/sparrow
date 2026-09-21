#!/usr/bin/env bash
# Build the browser WASM host adapter into ./dist.
set -euo pipefail
cd "$(dirname "$0")"

mkdir -p dist
OUT="${OUT:-dist/sparrow.wasm}"

echo "building $OUT (GOOS=js GOARCH=wasm)..."
GOOS=js GOARCH=wasm go build -trimpath -ldflags="-s -w" -o "$OUT" .

GOROOT="$(go env GOROOT)"
cp "$GOROOT/lib/wasm/wasm_exec.js" dist/wasm_exec.js 2>/dev/null \
  || cp "$GOROOT/misc/wasm/wasm_exec.js" dist/wasm_exec.js

BYTES="$(wc -c < "$OUT" | tr -d ' ')"
echo "wrote $OUT ($BYTES bytes)"
echo "copied dist/wasm_exec.js"
if command -v gzip >/dev/null 2>&1; then
  gzip -kf "$OUT"
  GBYTES="$(wc -c < "${OUT}.gz" | tr -d ' ')"
  echo "wrote ${OUT}.gz ($GBYTES bytes)"
fi
