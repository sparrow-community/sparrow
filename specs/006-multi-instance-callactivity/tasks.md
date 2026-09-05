---
description: "Task list for Multi-Instance Call Activity"
---

# Tasks: Multi-Instance Call Activity

**Input**: Design documents from `/specs/006-multi-instance-callactivity/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Included — constitution requires Event-intent tests; spec SC-001–SC-007.

**Organization**: Tasks are grouped by user story so each story can be implemented and tested independently.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1–US5
- Paths are repository-relative

## Phase 1: Setup

**Purpose**: Baseline before changing Call Activity / MI wiring

- [x] T001 Run baseline: `go test ./processing/ -run 'CallActivity|CrossDeploy|MultiInstance' -count=1`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Deploy + publication plumbing every story needs; no user-visible MI Call Activity until US1

**⚠️ CRITICAL**: No user story work until this phase completes

- [x] T002 Stop rejecting multi-instance Call Activity in `validateCallActivity` in `processing/deploy/call_activity.go`
- [x] T003 Index Call Activity loop characteristics via `indexMultiInstance` in `processing/deploy/deploy.go` (CallActivities path / `indexCallActivities`)
- [x] T004 Add optional snapped child input variables on `handlers.Publication` in `processing/handlers/handler.go`
- [x] T005 Propagate publications from `enterWithLoopIndex` in `runMultiInstanceStart` in `processing/multi_instance_executor.go` (and wire returned pubs through `Executor.Enter` in `processing/executor.go`)
- [x] T006 Prefer snapped `Publication` inputs in `startCalledInstance` in `processing/call_child.go`
- [x] T007 On MI cancel, queue `PublicationTerminateChild` for each inner Call Activity `CalledProcessInstanceID` in `processing/multi_instance_executor.go` (flush via existing publication path)

**Checkpoint**: Existing non-MI Call Activity and MI User/Service/SubProcess tests still pass

---

## Phase 3: User Story 1 — Parallel multi-instance Call Activity (P1) 🎯 MVP

**Goal**: Parallel MI Call Activity with fixed cardinality; N distinct children; host completes once after all children

**Independent Test**: Quickstart Story 1 — cardinality 3 → three children; complete all → parent ends once

### Tests for User Story 1

- [x] T008 [P] [US1] Add `processing/testdata/m10_mi_call_parallel.bpmn`
- [x] T009 [US1] Add parallel MI Call Activity join test in `processing/call_activity_test.go` (or `processing/multi_instance_test.go`)
- [x] T010 [US1] Add deploy-accepts-MI-CallActivity test in `processing/call_activity_test.go`

### Implementation for User Story 1

- [x] T011 [US1] Branch `CallActivityHandler.OnEnter` for MI host (`multiInstanceHostEnter`) vs inner (start child) vs non-MI in `processing/handlers/call_activity.go`
- [x] T012 [US1] Branch `CallActivityHandler.OnComplete` for MI inner complete vs host complete (boundaries + compensation on host) in `processing/handlers/call_activity.go`
- [x] T013 [US1] Snapshot mapped inputs onto `PublicationStartChild` at inner enter in `processing/handlers/call_activity.go`
- [x] T014 [US1] Treat zero/negative cardinality as immediate host complete (reuse MI host path) for Call Activity in `processing/handlers/call_activity.go`

**Checkpoint**: Story 1 quickstart passes; `m4_call_*` / `m8_cross_call_*` / existing MI suites still pass

---

## Phase 4: User Story 2 — Sequential multi-instance Call Activity (P1)

**Goal**: Sequential MI Call Activity; at most one active child at a time

**Independent Test**: Quickstart Story 2 — cardinality 3; one child at a time; three ordered completes

### Tests for User Story 2

- [x] T015 [P] [US2] Add `processing/testdata/m10_mi_call_sequential.bpmn`
- [x] T016 [US2] Add sequential-only-one-active-child test in `processing/call_activity_test.go`

### Implementation for User Story 2

- [x] T017 [US2] Verify sequential next-index spawn after inner Call Activity complete reuses existing MI executor path in `processing/multi_instance_executor.go` / `processing/handlers/call_activity.go` (fix gaps only)

**Checkpoint**: Story 1 still passes; Story 2 sequential fixture passes

---

## Phase 5: User Story 3 — Collection input and IO mapping (P2)

**Goal**: Collection-driven child count; per-child mapped inputs; empty collection immediate complete; optional output collection

**Independent Test**: Quickstart Story 3 — collection size N; empty collection; mapped vars on children

### Tests for User Story 3

- [x] T018 [P] [US3] Add `processing/testdata/m10_mi_call_collection.bpmn` (and empty-collection variant if split)
- [x] T019 [US3] Add collection input / empty collection / optional output collection tests in `processing/call_activity_test.go`

### Implementation for User Story 3

- [x] T020 [US3] Ensure MI collection element vars + Call Activity input mappings are snapped per index at StartChild in `processing/handlers/call_activity.go` / `processing/call_child.go`
- [x] T021 [US3] Assemble Call Activity output mappings into MI output collection on host complete in `processing/multi_instance_executor.go` / `processing/handlers/call_activity.go`

**Checkpoint**: Stories 1–2 still pass; Story 3 collection fixtures pass

---

## Phase 6: User Story 4 — Completion condition and interrupting boundary (P2)

**Goal**: Early join cancels remaining children; interrupting boundary terminates all children and takes boundary path once

**Independent Test**: Quickstart Story 4 — 2-of-5 completion; timer boundary on MI Call Activity

### Tests for User Story 4

- [x] T022 [P] [US4] Add `processing/testdata/m10_mi_call_completion.bpmn`
- [x] T023 [P] [US4] Add `processing/testdata/m10_mi_call_timer_boundary.bpmn`
- [x] T024 [US4] Add completion-condition early join + child terminate test in `processing/call_activity_test.go`
- [x] T025 [US4] Add interrupting timer boundary on MI Call Activity test in `processing/call_activity_test.go` or `processing/boundary_test.go`

### Implementation for User Story 4

- [x] T026 [US4] Confirm completion-condition cancel path terminates Call Activity children (T007) end-to-end in `processing/multi_instance_executor.go`
- [x] T027 [US4] Confirm MI host boundary cancel (`MultiInstanceCancel`) terminates all Call Activity children in `processing/handlers/boundary_event.go` / executor cancel path

**Checkpoint**: Stories 1–3 still pass; Story 4 fixtures pass

---

## Phase 7: User Story 5 — Recover mid multi-instance call (P2)

**Goal**: Recover mid-loop preserves counters, inners, and children; finishing remaining work matches no-crash outcome

**Independent Test**: Quickstart Story 5 — 1 of 3 children done; Recover; complete rest

### Tests for User Story 5

- [x] T028 [US5] Add Recover mid-loop parallel MI Call Activity test in `processing/recover_test.go`
- [x] T029 [P] [US5] Add cross-deploy Recover mid-loop fixture/test variant in `processing/recover_test.go` (reuse `m8` callee or `m10_mi_call_cross_*.bpmn`)

### Implementation for User Story 5

- [x] T030 [US5] Fix any Recover gaps for MI Call Activity host/inner/`CalledDeploymentID` in `processing/recover.go` / `processing/call_child.go` / projection (only if tests fail)

**Checkpoint**: Story 5 Recover tests pass; prior stories still pass

---

## Phase 8: Polish & Cross-Cutting

- [x] T031 Run full suite: `go test ./processing/ ./gateway/ ./protocol/proto/event/v1/`
- [x] T032 Update `AGENTS.md`: move Multi-instance Call Activity to Implemented; clear Active increment

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Immediate
- **Foundational (Phase 2)**: After Setup — **blocks all stories**
- **US1 (Phase 3)**: After Foundational — MVP
- **US2 (Phase 4)**: After US1 handler/MI spawn works (sequential reuses same enter/complete)
- **US3 (Phase 5)**: After US1 (needs snapped inputs)
- **US4 (Phase 6)**: After Foundational T007 + US1 (cancel + host boundaries)
- **US5 (Phase 7)**: After US1 (mid-loop state); ideally after US3/US4 if covering those paths
- **Polish (Phase 8)**: After desired stories complete

### User Story Dependencies

- **US1**: No story deps after Foundational
- **US2**: Depends on US1 Call Activity MI enter/complete
- **US3**: Depends on US1 + T004/T006 snapshot path
- **US4**: Depends on US1 + T007 terminate-all
- **US5**: Depends on US1; cross-deploy variant depends on `004` behavior unchanged

### Parallel Opportunities

- T008 fixtures parallel with early US1 impl after Foundational
- T015 / T018 / T022 / T023 fixture files in parallel across stories once Foundational done
- T028 and T029 Recover tests can be authored in parallel

---

## Parallel Example: User Story 1

```text
# After Foundational:
T008 Add m10_mi_call_parallel.bpmn
T011–T014 CallActivityHandler MI host/inner wiring
T009–T010 Tests once handler starts children
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Phase 1 Setup baseline
2. Phase 2 Foundational (deploy + pubs + terminate children)
3. Phase 3 US1 parallel MI Call Activity
4. **STOP and VALIDATE** quickstart Story 1 + regression

### Incremental Delivery

1. US2 sequential → validate
2. US3 collection/IO → validate
3. US4 completion + boundary → validate
4. US5 Recover → validate
5. Polish + `AGENTS.md`

---

## Notes

- No proto changes expected.
- Critical Foundational: do not ship MI Call Activity without T005 (pubs from MI start) or children never start.
- Compensation subscribe only on successful **host** complete (existing Call Activity OnComplete host path).
- Instantiating event-based gateway and live migration remain out of scope.
