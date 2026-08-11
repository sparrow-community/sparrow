#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")"
buf lint
buf generate
