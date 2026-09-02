# Implementation Plan: Call Activity boundary and compensation parity

**Branch**: `005-callactivity-boundary-compensation` | **Date**: 2026-09-02 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/005-callactivity-boundary-compensation/spec.md`

## Summary

Bring Call Activity to **parity with SubProcess** for deploy-time boundary attachment and runtime boundary/compensation lifecycle. Today `validateBoundaryHost` rejects boundaries on `callActivity`; `CallActivityHandler` does not arm timer/message/signal boundaries on ACTIVATED or subscribe compensation on COMPLETED — while `BoundaryEventHandler` and `fireActivityErrorBoundary` already terminate child instances when interrupting paths fire. This increment is mostly **deploy + handler wiring** plus tests and Recover coverage; no new ledger subjects or RPCs.

## Technical Context

**Language/Version**: Go 1.26.5 (workspace `go.work`)

**Primary Dependencies**: in-tree `bpmn`, `protocol` (existing `ActivityPayload` boundary fields), `processing` (`deploy/deploy.go`, `handlers/call_activity.go`, `handlers/boundary_event.go`, `errors.go`, `projection`, `recover`)

**Storage**: append-only `EventLog` + `deploy.Store`; boundary/compensation state on existing token projection fields

**Testing**: `go test ./processing/ ./gateway/ ./protocol/proto/event/v1/`; new fixtures `m9_call_boundary_*.bpmn`, `m9_call_compensate_*.bpmn`; regression on `m4_call_*`, `m8_cross_call_*`, `boundary_test.go`

**Target Platform**: single-node process (`gateway/cmd/sparrow`)

**Project Type**: Go multi-module BPMN execution engine

**Performance Goals**: correctness and recoverability; reuse existing boundary attach/disarm paths

**Constraints**: constitution (Element subject unchanged, no parallel graph); one boundary per kind per host (unchanged); no proto changes expected

**Scale/Scope**: four user stories (interrupting boundaries, non-interrupting, compensation, Recover). Deferred: multi-instance Call Activity boundaries, multiple same-kind boundaries, live migration.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status |
|-----------|--------|
| I. Engine, not product suite | Pass — handler/deploy semantics only |
| II. Ledger is source of truth | Pass — boundary/compensation facts remain existing BOUNDARY_EVENT / host intents; Recover replays token fields |
| III. Element is the subject | Pass — no new Element.Type; CALL_ACTIVITY + BOUNDARY_EVENT intents |
| IV. No parallel graph; thin modules | Pass — extend `validateBoundaryHost`, `CallActivityHandler`, reuse `attachBoundary` / `subscribeCompensation` |
| V. Standard BPMN | Pass — standard boundaryEvent attachment to callActivity |
| VI. Serial commands, explicit rejection | Pass — unchanged; boundary Complete uses existing host token lock |

Post-design re-check: still pass. Interrupting paths reuse `PublicationTerminateChild` and existing `terminateChildPub` in error boundary path.

## Project Structure

### Documentation (this feature)

```text
specs/005-callactivity-boundary-compensation/
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
├── deploy/deploy.go           # validateBoundaryHost: add callActivity
├── handlers/
│   ├── call_activity.go       # OnEnter: attachBoundary; OnComplete: cancel + subscribeCompensation
│   └── boundary_event.go      # (verify) interrupting CALL_ACTIVITY → TerminateChild — already present
├── errors.go                  # fireActivityErrorBoundary + terminateChildPub — already present
├── call_child.go              # terminateCalledInstance — unchanged
├── boundary_test.go           # extend or parallel call-activity boundary cases
├── call_activity_test.go      # interrupting timer/message on call
├── compensation_test.go       # call activity compensation fixture
└── recover_test.go            # mid-call boundary + compensation subscription

processing/testdata/
└── m9_call_*.bpmn
```

**Structure Decision**: Minimal diff in `deploy` + `call_activity` handler; leverage shared `attachBoundary`, `cancelAttachedBoundary`, `subscribeCompensation` from `boundary_event.go`.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | — | — |

## Phase 0 & 1 Outputs

- [research.md](./research.md) — gap analysis vs SubProcess; reuse existing interrupt/error child terminate
- [data-model.md](./data-model.md) — host token boundary fields, compensation subscription, deploy host set
- [contracts/engine.md](./contracts/engine.md) — deploy acceptance rules; no new RPC
- [quickstart.md](./quickstart.md) — per-story validation steps

## Implementation Notes (for tasks phase)

1. **Deploy**: Add `callActivities` loop in `validateBoundaryHost` (mirror subProcess entry in `deploy.go`).
2. **OnEnter**: After child id allocation, call `attachBoundary(dep, elementID, now, payload)` and embed result in CALL_ACTIVITY ACTIVATED `ActivityPayload` (same pattern as `user_task.go`).
3. **OnComplete**: Append `cancelAttachedBoundary` + `subscribeCompensation` before `TakeOutgoing` (mirror `sub_process.go` OnComplete for non-MI host).
4. **Interrupting runtime**: Verify `BoundaryEventHandler.OnComplete` `PublicationTerminateChild` for `TYPE_CALL_ACTIVITY` — already implemented; add tests only unless gap found.
5. **Error boundary**: `fireActivityErrorBoundary` already calls `terminateChildPub` — add deploy + arm + test with Call Activity host.
6. **Non-interrupting**: Reuse existing `BoundaryEventHandler` non-interrupting branch; ensure host token stays `waiting` with child id after spawn.
7. **Compensation**: `subscribeCompensation` on Call Activity COMPLETED; compensate throw uses existing scope queue — test with caller-side handler.
8. **Recover**: Replay restores `DueUnixMs`, boundary ids, `CalledProcessInstanceID`, `CompensationSubs` — extend `recover_test.go`.
9. **Fixtures**: `m9_call_timer_boundary.bpmn`, `m9_call_message_non_interrupt.bpmn`, `m9_call_compensate.bpmn`, optional cross-deploy interrupt variant.
10. **AGENTS.md**: move item from Remaining to Implemented after implement phase.
