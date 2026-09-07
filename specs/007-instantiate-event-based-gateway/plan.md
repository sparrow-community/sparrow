# Implementation Plan: Instantiate Event-Based Gateway

**Branch**: `007-instantiate-event-based-gateway` | **Date**: 2026-09-07 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/007-instantiate-event-based-gateway/spec.md`

## Summary

Allow an **exclusive** event-based gateway with `instantiate="true"` to be the process entry when there is no none startEvent. Deploy stops rejecting instantiate; validation requires no incoming flows, ≥2 catch targets, exclusive type, and no startEvent+instantiate mix. `StartEventID` / CreateInstance / Recover enter that gateway element so existing exclusive fork + sibling-cancel behavior arms the start race. Parallel instantiate and message-only instance creation stay out of scope.

## Technical Context

**Language/Version**: Go 1.26.5 (workspace `go.work`)

**Primary Dependencies**: in-tree `bpmn` (`EventBasedGateway.Instantiate`), `processing` (`deploy/deploy.go` validate + `StartEventID`, `handlers/event_based_gateway.go`, `engine.go` CreateInstance, `recover.go` redrive)

**Storage**: append-only `EventLog` + `deploy.Store`; no new ledger subjects

**Testing**: `go test ./processing/ ./gateway/ ./protocol/proto/event/v1/`; fixtures `m11_instantiate_ebg_*.bpmn`; regression mid-process EBG (`m2_*`, `m3_*`)

**Target Platform**: single-node process (`gateway/cmd/sparrow`)

**Project Type**: Go multi-module BPMN execution engine

**Performance Goals**: correctness and recoverability; reuse exclusive EBG cancel

**Constraints**: constitution (no new Element.Type); no proto changes; CreateInstance remains the only instance-creation API (FR-010)

**Scale/Scope**: three user stories (happy path, invalid deploy, Recover). Out of scope: parallel instantiate, subscription-start without CreateInstance, live migration.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status |
|-----------|--------|
| I. Engine, not product suite | Pass — deploy + entry resolve + existing handler |
| II. Ledger is source of truth | Pass — same EVENT_BASED_GATEWAY + catch events; Recover redrives CreateInstance into entry |
| III. Element is the subject | Pass — existing `TYPE_EVENT_BASED_GATEWAY` / catch types |
| IV. No parallel graph; thin modules | Pass — deploy validation + `StartEventID` semantics; no new package |
| V. Standard BPMN | Pass — `instantiate` on eventBasedGateway |
| VI. Serial commands, explicit rejection | Pass — invalid shapes fail deploy; CreateInstance still serial per instance |

Post-design re-check: still pass.

## Project Structure

### Documentation (this feature)

```text
specs/007-instantiate-event-based-gateway/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── engine.md
└── tasks.md              # Phase 2 (/speckit-tasks)
```

### Source Code (repository root)

```text
processing/
├── deploy/deploy.go              # validate instantiate; entry via StartEventID
├── engine.go / recover.go        # CreateInstance / redrive use entry id (unchanged call sites if StartEventID extended)
├── handlers/event_based_gateway.go  # unchanged behavior expected
├── event_based_gateway_test.go / recover_test.go
└── testdata/m11_instantiate_ebg_*.bpmn
```

**Structure Decision**: Minimal deploy + process-entry resolution; runtime fork/cancel already works for exclusive EBG.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | — | — |
