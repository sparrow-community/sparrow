# Implementation Plan: Compatibility Baseline

**Branch**: `036-compatibility-baseline` | **Date**: 2026-10-16 | **Spec**: [spec.md](./spec.md)

## Summary

Deliver constitution Principle VII’s first enforcement layer without changing BPMN semantics: publish `COMPATIBILITY.md` (versioning + surface inventory), commit a golden `Open` data directory and a replay test, freeze today’s `.proto` tree as `protocol/breaking-baseline/` for `buf breaking`, and wire both checks into CI. Bootstrap the proto baseline because no protocol release tag exists yet.

## Technical Context

**Language/Version**: Go 1.26.5 · buf CLI in CI  
**Primary Dependencies**: `event.v1` length-delimited EventLog · `processing.Open`  
**Storage**: file EventLog + DirStore (golden under `processing/testdata/compat/`)  
**Testing**: `go test` in `processing`; `buf breaking` in CI  
**Project Type**: library / kernel  
**Constraints**: no product UI; additive-only wire; keep Apache-2.0 headers  

## Constitution Check

- Principle I: docs + tests + CI only; no product suite.  
- Principle II/VII: golden replay is the ledger promise.  
- Module boundaries: no new top-level modules.  
- Quality gates: golden replay + `buf breaking` satisfy the new VII gate for this increment.

## Project Structure

```text
COMPATIBILITY.md
README.md                          # point at COMPATIBILITY.md
AGENTS.md / AGENTS.zh.md           # Planned P1 status
.github/workflows/ci.yml           # buf breaking job
protocol/breaking-baseline/        # bootstrap .proto snapshot
protocol/proto/build.sh            # optional local breaking check note
processing/testdata/compat/baseline/
processing/compat_golden_test.go
processing/scripts/regen-compat-baseline.sh
specs/036-compatibility-baseline/
```

## Design notes

1. **Versioning**: stay on WASM calendar alpha until 1.0 criteria in COMPATIBILITY.md are met; Go modules remain untagged (pseudo-version) until an explicit `v0.x` / `v1.0` tag plan is cut.  
2. **Golden**: generate once via script into `testdata/compat/baseline/{events.log,deployments/,expected.json}`; test asserts Open → status/waiting/variables.  
3. **buf**: CI `buf breaking --against ../breaking-baseline` from `protocol/proto`. Document switch to `.git#tag=…` when a protocol tag exists.  
4. **Out of scope**: SQLite store, Node/Deno WASM (P2); MCP adapter (P3 consumer).
