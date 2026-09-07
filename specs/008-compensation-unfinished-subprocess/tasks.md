---
description: "Task list for Compensation into Unfinished SubProcess"
---

# Tasks: Compensation into Unfinished SubProcess

**Input**: Design documents from `/specs/008-compensation-unfinished-subprocess/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Included — SC-001–SC-005.

## Format: `[ID] [P?] [Story] Description`

## Phase 1: Setup

- [x] T001 Run baseline: `go test ./processing/ -run 'Compensation|Compensate' -count=1`

## Phase 2: Foundational

- [x] T002 Add helpers to detect unfinished child SubProcesses and `activityInScopeTree` in `processing/executor.go` (or small helper file)
- [x] T003 Extend `startCompensation` to terminate unfinished SPs and merge under-SP subscriptions into the handler queue in `processing/executor.go`

**Checkpoint**: Existing compensation tests still pass

## Phase 3: User Story 1 — Broadcast unfinished SP (P1) 🎯 MVP

- [x] T004 [P] [US1] Add `processing/testdata/m12_compensate_unfinished_sp_parallel.bpmn`
- [x] T005 [US1] Add broadcast unfinished-SP compensation test in `processing/compensation_test.go`
- [x] T006 [US1] Add empty-inner-subscription unfinished SP test (terminate only) if not covered by T005

## Phase 4: User Story 2 — Targeted activityRef (P1)

- [x] T007 [P] [US2] Add `processing/testdata/m12_compensate_unfinished_sp_targeted.bpmn`
- [x] T008 [US2] Add targeted unfinished-SP test (sibling subscription preserved) in `processing/compensation_test.go`

## Phase 5: User Story 3 — Recover (P2)

- [x] T009 [US3] Add Recover mid unfinished-SP compensation test in `processing/compensation_test.go`
- [x] T009b Fix projection `TERMINATED` to drop tokens on replay (`processing/projection/instance.go`) so Recover does not revive cancelled scope work

## Phase 6: Polish

- [x] T010 Update `AGENTS.md` for 008 shipped
- [x] T011 Run `go test ./processing/ ./gateway/ ./protocol/proto/event/v1/`

## Dependencies

Phase 1 → 2 → US1 → US2 → US3 → Polish
