# Implementation Plan: Compensation into Unfinished SubProcess

**Branch**: `008-compensation-unfinished-subprocess` | **Date**: 2026-09-07 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/008-compensation-unfinished-subprocess/spec.md`

## Summary

Extend `startCompensation` so a parent-scope compensate throw can **enter unfinished embedded SubProcesses**: terminate active tokens in that SubProcess (reuse `terminateScope`), then queue compensation handlers for subscriptions whose activities live under that SubProcess (reverse Seq order), merged with existing same-scope completed subscriptions for broadcast. Targeted `activityRef` to an unfinished SubProcess enters only that SubProcess. Recover continues to rebuild PendingCompensation from waiting throw + handler.

## Technical Context

**Language/Version**: Go 1.26.5 (workspace `go.work`)

**Primary Dependencies**: `processing/executor.go` (`startCompensation` / `advanceCompensation`), `processing/scope_terminate.go` / `terminateScope`, `processing/deploy` (`ScopeOf`, `CompensationByBoundary`, `CompensateActivityRef`), existing compensation projection

**Storage**: append-only EventLog; existing `CompensationSubs` + `PendingCompensation` (no new ledger subjects)

**Testing**: `go test ./processing/ -run 'Compensation|Compensate'`; fixtures `m12_compensate_unfinished_sp_*.bpmn`; regression `m3_compensation`, `m4_compensate_subprocess`, `m9_call_compensate`

**Target Platform**: single-node (`gateway/cmd/sparrow`)

**Project Type**: Go multi-module BPMN execution engine

**Performance Goals**: correctness + Recover; reuse terminateScope

**Constraints**: constitution (no new Element.Type); no proto changes; Call Activity / Transaction / cancel / compensation ESP out of scope

**Scale/Scope**: three user stories (broadcast unfinished, targeted activityRef, Recover)

## Constitution Check

| Principle | Status |
|-----------|--------|
| I. Engine, not product suite | Pass |
| II. Ledger is source of truth | Pass — terminate + boundary COMPLETED + handler enter recorded as today |
| III. Element is the subject | Pass — existing types |
| IV. Thin modules | Pass — executor compensation path + terminate helpers |
| V. Standard BPMN | Pass — compensate throw into active embedded SubProcess completeness |
| VI. Serial commands, explicit rejection | Pass |

Post-design re-check: still pass.

## Project Structure

### Documentation (this feature)

```text
specs/008-compensation-unfinished-subprocess/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/engine.md
└── tasks.md
```

### Source Code (repository root)

```text
processing/
├── executor.go                 # startCompensation unfinished-SP entry
├── scope_terminate.go          # reuse terminateScope
├── compensation_test.go
└── testdata/m12_compensate_unfinished_sp_*.bpmn
```

**Structure Decision**: Extend compensation start only; no new packages.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | — | — |
