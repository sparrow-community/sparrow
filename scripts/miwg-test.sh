#!/usr/bin/env bash
# Run BPMN MIWG fixture load + kernel deploy/run suite (CI-equivalent).
set -euo pipefail
root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"

echo "==> bpmn TestMIWGFixturesLoad (pin $(cat bpmn/test/MIWG_UPSTREAM_SHA))"
go test -C bpmn -run TestMIWGFixturesLoad ./...

echo "==> processing TestMIWGKernelFixtures"
go test -C processing -run TestMIWGKernelFixtures ./...

echo "MIWG suite OK"
