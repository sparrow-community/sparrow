# Implementation Plan: Multi-Instance Call Activity

**Branch**: `006-multi-instance-callactivity` | **Date**: 2026-09-05 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/006-multi-instance-callactivity/spec.md`

## Summary

Extend existing multi-instance loop machinery (`002`) to **Call Activity**: each inner iteration starts one called process instance (same-file or cross-deployment). Deploy today rejects MI on Call Activity; runtime `CallActivityHandler` always starts a single child. This increment removes that reject, indexes MI specs on Call Activity, wires host/inner enter-complete like User Task / SubProcess, ensures `PublicationStartChild` pubs from MI spawn are not dropped, snapshots per-iteration inputs for deferred child start, and terminates all active children on early completion or interrupting boundary.

## Technical Context

**Language/Version**: Go 1.26.5 (workspace `go.work`)

**Primary Dependencies**: in-tree `bpmn` (Call Activity already embeds loop characteristics), `protocol` (existing `ActivityPayload.loop_*` / `called_process_instance_id`), `processing` (`deploy/call_activity.go`, `deploy/multi_instance.go`, `handlers/call_activity.go`, `handlers/multi_instance.go`, `multi_instance_executor.go`, `call_child.go`, `projection`)

**Storage**: append-only `EventLog` + `deploy.Store`; loop state in existing `MultiInstanceLoops` projection; child linkage on inner tokens

**Testing**: `go test ./processing/ ./gateway/ ./protocol/proto/event/v1/`; fixtures `m10_mi_call_*.bpmn`; regression `CallActivity|CrossDeploy|MultiInstance`

**Target Platform**: single-node process (`gateway/cmd/sparrow`)

**Project Type**: Go multi-module BPMN execution engine

**Performance Goals**: correctness and recoverability; reuse MI join/cancel and Call Activity start/resume/terminate

**Constraints**: constitution (no new Element.Type); no proto changes expected; one COMMAND per parent instance lock; child start/terminate still via deferred publications

**Scale/Scope**: five user stories (parallel, sequential, collection+IO, completion/boundary, Recover). Out of scope: instantiate event-based gateway, live migration, complex behavior defs.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status |
|-----------|--------|
| I. Engine, not product suite | Pass — handler/deploy/executor only |
| II. Ledger is source of truth | Pass — existing CALL_ACTIVITY + loop payload fields; Recover rebuilds loops and child links |
| III. Element is the subject | Pass — same `TYPE_CALL_ACTIVITY`; inners distinguished by token + `loop_instance_index` |
| IV. No parallel graph; thin modules | Pass — extend CallActivity handler + MI executor + call_child; no new package |
| V. Standard BPMN | Pass — `multiInstanceLoopCharacteristics` on callActivity |
| VI. Serial commands, explicit rejection | Pass — parent Completes via resume; child lifecycle on child instance lock |

Post-design re-check: still pass. Critical fix: MI spawn must return/flush `PublicationStartChild`; cancel path must terminate children.

## Project Structure

### Documentation (this feature)

```text
specs/006-multi-instance-callactivity/
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
├── deploy/call_activity.go          # stop rejecting MI; keep CallActivity spec
├── deploy/deploy.go                 # indexMultiInstance for CallActivities
├── handlers/call_activity.go        # MI host + inner enter/complete paths
├── handlers/multi_instance.go       # reuse host enter / inner complete helpers
├── handlers/handler.go              # Publication: optional snapped child inputs
├── multi_instance_executor.go       # return pubs from MI start; terminate children on cancel
├── call_child.go                    # use snapped inputs; terminate-all helper if needed
├── multi_instance_test.go / call_activity_test.go / recover_test.go
└── testdata/m10_mi_call_*.bpmn
```

**Structure Decision**: Minimal extension of `002` + `004`/`005` paths; no new ledger subjects or RPCs.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | — | — |
