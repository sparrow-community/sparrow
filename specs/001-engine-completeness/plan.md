# Implementation Plan: Engine Completeness (next increment)

**Branch**: `001-engine-completeness` | **Date**: 2026-08-26 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/001-engine-completeness/spec.md`

## Summary

Close the next executable-BPMN gaps without leaving the four-module layout: Call Activity becomes a real child process instance (with IO mapping), error Event Sub-Process is armed like the existing message/timer/signal starts, and deployments of the same process id keep numbered revisions that new starts can select. Cross-instance resume reuses the existing post-lock publication pattern; the ledger subject remains Element.

## Technical Context

**Language/Version**: Go 1.26.5 (workspace `go.work`)

**Primary Dependencies**: in-tree `bpmn`, `protocol` (Protobuf + gRPC generated Go), `processing` handlers/executor, `gateway` EngineService/JobService

**Storage**: append-only `log.EventLog` + `deploy.Store` + optional `runtime.Store` (Memory/File in-tree)

**Testing**: `go test` (`processing` fixture tests are the semantic contract; `gateway` RPC tests; `protocol/proto/event/v1`)

**Target Platform**: single-node process (`gateway/cmd/sparrow`)

**Project Type**: Go multi-module engine + gRPC adapter

**Performance Goals**: correctness and recoverability over throughput; one COMMAND at a time per `process_instance_id`

**Constraints**: constitution (ledger subject, no parallel graph, no gRPC in `processing`, no product-suite scope); `buf` FILE compatibility for proto changes

**Scale/Scope**: one feature directory covering three sequential stories (Call Activity instance+IO, error Event Sub-Process, revision coexistence). Deferred items stay out of this plan.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status |
|-----------|--------|
| I. Engine, not product suite | Pass — no UI/ops surfaces |
| II. Ledger is source of truth | Pass — parent/child link and mappings live on PROCESS/CALL_ACTIVITY Event payloads; leases/buffers unchanged |
| III. Element is the subject | Pass — no Job/Timer/Message subjects; child instance is still PROCESS + FlowElements |
| IV. No parallel graph; thin modules | Pass — still `element.Process`; gRPC stays in `gateway`; proto in `protocol` |
| V. Standard BPMN | Pass — Call Activity, IO associations, error start Event Sub-Process, process versioning |
| VI. Serial commands, explicit rejection | Pass — each instance keeps its own lock; caller resume is a later COMMAND on the parent id after the child lock is released (same pattern as message/signal Publication) |

Post-design re-check: still pass. Cross-instance resume is an Effect/publication, not a second ledger ontology.

## Project Structure

### Documentation (this feature)

```text
specs/001-engine-completeness/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── engine.md
└── tasks.md
```

### Source Code (repository root)

```text
bpmn/                     # XML model (IO associations already on Activity)
protocol/proto/           # event.v1 payloads + engine.v1 query/start fields
protocol/gen/go/          # generated; via protocol/proto/build.sh
processing/
├── deploy/               # CallActivity spec, mappings, version index, error Event Sub-Process
├── handlers/             # call_activity.go, start_event.go, end_event.go
├── projection/           # parent/child fields, called instance id on host token
├── executor.go           # spawn child instance; resume parent after child settle
└── *_test.go / testdata/ # fixtures per story
gateway/                  # map new EngineService fields
```

**Structure Decision**: Extend the existing four modules. Do not add packages for "child engine" or a parallel call graph.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | — | — |
