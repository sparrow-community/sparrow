#!/usr/bin/env bash
# Regenerate processing/testdata/compat/baseline after an intentional, reviewed
# additive change. Do not use this to silence a Principle VII failure.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
OUT="$ROOT/processing/testdata/compat/baseline"
rm -rf "$OUT"
mkdir -p "$OUT"

cd "$ROOT"
go run -C processing ./scripts/gencompatbaseline -out "$OUT"
echo "wrote $OUT"
