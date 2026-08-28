# Implementation Plan: Multi-Instance Activities

**Branch**: `002-multi-instance` | **Date**: 2026-08-28 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/002-multi-instance/spec.md`

## Summary

Add BPMN multi-instance loop execution for User Task, Service Task, and embedded Sub-Process: parallel and sequential modes, fixed cardinality and collection input, default-all and expression-based completion conditions, optional output collection, and Recover-safe loop counters. Loop host remains the same FlowElement subject; inner iterations are distinguished by token/payload loop index fields and projection loop state — no new ledger ontology.

## Technical Context

**Language/Version**: Go 1.26.5 (workspace `go.work`)

**Primary Dependencies**: in-tree `bpmn` (`element.MultiInstanceLoopCharacteristics` already parsed), `protocol` (Protobuf + generated Go), `processing` handlers/executor/deploy, `processing/expr` for completion conditions

**Storage**: append-only `log.EventLog` + `deploy.Store` + optional `runtime.Store` (unchanged)

**Testing**: `go test` (`processing` fixture tests per story; `gateway` if proto/query fields added; `protocol/proto/event/v1`)

**Target Platform**: single-node process (`gateway/cmd/sparrow`)

**Project Type**: Go multi-module engine + gRPC adapter

**Performance Goals**: correctness and recoverability over throughput; one COMMAND at a time per `process_instance_id`

**Constraints**: constitution (Element subject, no parallel graph model, no gRPC in `processing`, no product-suite scope); `buf` FILE compatibility for proto changes

**Scale/Scope**: five user stories (parallel cardinality, sequential, collection input, completion condition, MI sub-process). Deferred: MI Call Activity, complex behavior, none/one behavior events.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status |
|-----------|--------|
| I. Engine, not product suite | Pass — no UI/ops surfaces |
| II. Ledger is source of truth | Pass — loop counters and inner-instance lifecycle recorded on existing ACTIVITY/PROCESS intents with additive payload fields; projection rebuilds loop state |
| III. Element is the subject | Pass — inner instances are still the same `Element.Type` (USER_TASK, etc.) with loop index on payload; no `LOOP_INSTANCE` subject |
| IV. No parallel graph; thin modules | Pass — compile MI spec into `deploy`; handler + executor extensions; proto in `protocol` |
| V. Standard BPMN | Pass — `multiInstanceLoopCharacteristics` from OMG model |
| VI. Serial commands, explicit rejection | Pass — all inner instances share one `process_instance_id` lock; serial COMMAND queue orders completes |

Post-design re-check: still pass. Loop state is projection + EVENT payload fields, not a second executable graph.

## Project Structure

### Documentation (this feature)

```text
specs/002-multi-instance/
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
bpmn/element/             # already has MultiInstanceLoopCharacteristics
protocol/proto/event/v1/  # additive ActivityPayload loop fields
protocol/gen/go/          # generated via protocol/proto/build.sh
processing/
├── deploy/               # compile MI spec (cardinality, collection, completion, sequential flag)
├── handlers/             # user_task, service_task, sub_process — MI enter/complete/cancel
├── projection/           # Token.loop_instance_index; MultiInstanceLoop projection state
├── executor.go           # spawn inner tokens, evaluate completion, cancel stragglers
├── expr/                 # reuse for completionCondition
└── *_test.go / testdata/ # m6_mi_*.bpmn fixtures per story
gateway/                  # map new Token/Complete fields if wire surface changes
```

**Structure Decision**: Extend the existing four modules. Do not add a `loop` package or parallel graph.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | — | — |
