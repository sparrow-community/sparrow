---
description: "Task list for Call Activity boundary and compensation parity"
---

# Tasks: Call Activity boundary and compensation parity

**Input**: Design documents from `/specs/005-callactivity-boundary-compensation/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Included — constitution requires Event-intent tests; spec SC-001–SC-006.

## Phase 1: Setup

- [x] T001 Run baseline: `go test ./processing/ -run 'CallActivity|CrossDeploy|Boundary' -count=1`

## Phase 2: Foundational

- [x] T002 Add `callActivity` to `validateBoundaryHost` in `processing/deploy/deploy.go`
- [x] T003 Arm boundaries in `CallActivityHandler.OnEnter` via `attachBoundary` in `processing/handlers/call_activity.go`
- [x] T004 Add `cancelAttachedBoundary` + `subscribeCompensation` in `CallActivityHandler.OnComplete` in `processing/handlers/call_activity.go`

## Phase 3: User Story 1 — Interrupting boundaries (P1) MVP

- [x] T005 [P] [US1] Add `processing/testdata/m9_call_timer_boundary.bpmn`
- [x] T006 [P] [US1] Add `processing/testdata/m9_call_error_boundary.bpmn`
- [x] T007 [US1] Add interrupting timer and error boundary tests in `processing/call_activity_test.go`
- [x] T008 [US1] Add deploy-accepts-call-activity-boundary test in `processing/boundary_test.go`

## Phase 4: User Story 2 — Non-interrupting boundaries (P1)

- [x] T009 [P] [US2] Add `processing/testdata/m9_call_message_non_interrupt.bpmn`
- [x] T010 [US2] Add non-interrupting message boundary test in `processing/call_activity_test.go`

## Phase 5: User Story 3 — Compensation (P2)

- [x] T011 [P] [US3] Add `processing/testdata/m9_call_compensate.bpmn`
- [x] T012 [US3] Add call activity compensation test in `processing/compensation_test.go`

## Phase 6: User Story 4 — Recover (P2)

- [x] T013 [US4] Add Recover mid-call timer boundary test in `processing/recover_test.go`
- [x] T014 [US4] Add Recover compensation handler preservation test in `processing/recover_test.go`

## Phase 7: Polish

- [x] T015 Run full suite: `go test ./processing/ ./gateway/ ./protocol/proto/event/v1/`
- [x] T016 Update `AGENTS.md`: move boundary/compensation parity to Implemented

## Notes

- `ThrowError` extended for `TYPE_CALL_ACTIVITY` and flushes `PublicationTerminateChild` (child termination on error boundary).
- `advanceCompensation` rebuilds `PendingCompensation` after replay when a compensate throw is in flight.
- Compensation boundary `ACTIVATING`/`COMPLETING` no longer corrupts host token `ElementID` on replay (`applyCompensationBoundaryToken`).
