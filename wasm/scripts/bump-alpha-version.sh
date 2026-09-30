#!/usr/bin/env bash
# Compute the next alpha version for @sparrow-community/wasm and write it to package.json.
#
# Scheme: YYYY.M.D-alpha.N  (UTC calendar date as major.minor.patch; prerelease alpha.N)
#   - YYYY / M / D are UTC year, month, day with no leading zeros (node-semver valid)
#   - N starts at 1 each UTC day and increments for each publish that day
#   - CI publishes with dist-tag "latest" (package.json publishConfig.tag) and
#     also promotes "alpha" to the same version when OIDC dist-tag is enabled
#   - No leading "v" in package.json (git tags may use vYYYY.M.D-alpha.N separately)
#
# Why not Jiaxing's display form v2026.09.28.1-alpha:
#   - npm/package.json rejects or strips a leading "v"
#   - semver numeric components cannot have leading zeros (09 invalid)
#   - four numeric segments before the prerelease is not MAJOR.MINOR.PATCH
# Mapping: v2026.09.28.1-alpha → 2026.9.28-alpha.1
#
# Queries the npm registry for already-published versions matching today's prefix.
# If the package (or today's prefix) does not exist yet, N=1.
set -euo pipefail
cd "$(dirname "$0")/.."

PKG_NAME="$(node -p "require('./package.json').name")"
YEAR="$(date -u +%Y)"
# %-m / %-d: no leading zeros (GNU date). Fallback for BSD date.
MONTH="$(date -u +%-m 2>/dev/null || date -u +%m | sed 's/^0//')"
DAY_NUM="$(date -u +%-d 2>/dev/null || date -u +%d | sed 's/^0//')"
BASE="${YEAR}.${MONTH}.${DAY_NUM}"
PREFIX="${BASE}-alpha."

max_n=0
# Prefer versions list; tolerate missing package (404 / empty).
if versions_json="$(npm view "$PKG_NAME" versions --json 2>/dev/null)"; then
  max_n="$(BASE="$BASE" node -e '
const raw = require("fs").readFileSync(0, "utf8").trim();
let versions = [];
try { versions = JSON.parse(raw); } catch { versions = []; }
if (!Array.isArray(versions)) versions = [versions].filter(Boolean);
const re = new RegExp("^" + process.env.BASE.replace(/\./g, "\\.") + "-alpha\\.(\\d+)$");
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

echo "alpha version → $VERSION (UTC ${BASE}, N=$next_n)"
