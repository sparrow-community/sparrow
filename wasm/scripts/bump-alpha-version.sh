#!/usr/bin/env bash
# Compute the next alpha version for @sparrow-community/wasm and write it to package.json.
#
# Scheme: 0.0.0-alpha.YYYYMMDD.N
#   - YYYYMMDD is UTC calendar day
#   - N starts at 1 each UTC day and increments for each publish that day
#   - Dist-tag remains "alpha" (see package.json publishConfig.tag)
#
# Queries the npm registry for already-published versions matching today's prefix.
# If the package (or today's prefix) does not exist yet, N=1.
set -euo pipefail
cd "$(dirname "$0")/.."

PKG_NAME="$(node -p "require('./package.json').name")"
DAY="$(date -u +%Y%m%d)"
PREFIX="0.0.0-alpha.${DAY}."

max_n=0
# Prefer versions list; tolerate missing package (404 / empty).
if versions_json="$(npm view "$PKG_NAME" versions --json 2>/dev/null)"; then
  max_n="$(DAY="$DAY" PREFIX="$PREFIX" node -e '
const raw = require("fs").readFileSync(0, "utf8").trim();
let versions = [];
try { versions = JSON.parse(raw); } catch { versions = []; }
if (!Array.isArray(versions)) versions = [versions].filter(Boolean);
const re = new RegExp("^0\\.0\\.0-alpha\\." + process.env.DAY + "\\.(\\d+)$");
let max = 0;
for (const v of versions) {
  const m = String(v).match(re);
  if (m) max = Math.max(max, Number(m[1]));
}
process.stdout.write(String(max));
' <<<"$versions_json")"
fi

next_n=$((max_n + 1))
VERSION="${PREFIX}${next_n}"

node -e '
const fs = require("fs");
const pkg = JSON.parse(fs.readFileSync("package.json", "utf8"));
pkg.version = process.argv[1];
fs.writeFileSync("package.json", JSON.stringify(pkg, null, 2) + "\n");
' "$VERSION"

echo "alpha version → $VERSION (UTC day $DAY, N=$next_n)"
