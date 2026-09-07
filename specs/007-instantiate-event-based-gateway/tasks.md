---
description: "Task list for Instantiate Event-Based Gateway"
---

# Tasks: Instantiate Event-Based Gateway

**Input**: Design documents from `/specs/007-instantiate-event-based-gateway/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Included — constitution requires Event-intent tests; spec SC-001–SC-005.

**Organization**: Tasks are grouped by user story so each story can be implemented and tested independently.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1–US3
- Paths are repository-relative

## Phase 1: Setup

**Purpose**: Baseline before changing EBG / entry resolve

- [x] T001 Run baseline: `go test ./processing/ -run 'EventBased' -count=1`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Deploy validation + process entry resolution every story needs

**⚠️ CRITICAL**: No user story work until this phase completes

- [x] T002 Allow instantiate exclusive EBG in `validateEventBasedGateways` in `processing/deploy/deploy.go` (reject Parallel instantiate; reject incoming; keep ≥2 catch outs; reject instantiate inside subProcess)
- [x] T003 Allow process with no startEvent when exactly one valid process-level instantiate exclusive EBG in `validateM1` / deploy path in `processing/deploy/deploy.go`
- [x] T004 Reject startEvent + any instantiate EBG coexistence in `processing/deploy/deploy.go`
- [x] T005 Extend `StartEventID` to return instantiate EBG id when no startEvent in `processing/deploy/deploy.go`

**Checkpoint**: Existing mid-process EBG tests still pass; invalid instantiate still fails deploy until US2 fixtures

---

## Phase 3: User Story 1 — Deploy and start via instantiate exclusive EBG (P1) 🎯 MVP

**Goal**: Instantiate-only process CreateInstance arms timer+message race; first event wins and cancels sibling

**Independent Test**: Quickstart Story 1

### Tests for User Story 1

- [x] T006 [P] [US1] Add `processing/testdata/m11_instantiate_ebg_timer_message.bpmn`
- [x] T007 [US1] Add deploy + CreateInstance arms two catches + timer-wins and message-wins tests in `processing/event_based_gateway_test.go`

### Implementation for User Story 1

- [x] T008 [US1] Verify CreateInstance / executor Enter path works with EBG entry (fix gaps only in `processing/engine.go` / start-event handler if needed)

**Checkpoint**: Story 1 passes; `m2_*` / `m3_*` EBG still pass

---

## Phase 4: User Story 2 — Reject invalid instantiate configurations (P1)

**Goal**: Fail-fast deploy for unsupported instantiate shapes

**Independent Test**: Quickstart Story 2

### Tests for User Story 2

- [x] T009 [P] [US2] Add invalid fixtures `m11_instantiate_ebg_incoming.bpmn`, `m11_instantiate_ebg_parallel.bpmn`, `m11_instantiate_ebg_with_start.bpmn`, `m11_instantiate_ebg_one_catch.bpmn` under `processing/testdata/`
- [x] T010 [US2] Add deploy-reject tests for each invalid fixture in `processing/event_based_gateway_test.go`

### Implementation for User Story 2

- [x] T011 [US2] Tighten error messages / edge coverage in `processing/deploy/deploy.go` if any reject case still deploys

**Checkpoint**: All invalid fixtures rejected; Story 1 still passes

---

## Phase 5: User Story 3 — Recover while instantiate race is armed (P2)

**Goal**: Recover mid-wait then message/timer win matches continuous run

**Independent Test**: Quickstart Story 3

### Tests for User Story 3

- [x] T012 [US3] Add Recover mid-instantiate-race test (message or timer win) in `processing/recover_test.go` or `processing/event_based_gateway_test.go`

### Implementation for User Story 3

- [x] T013 [US3] Confirm `redriveCreateInstance` uses extended `StartEventID` in `processing/recover.go` (fix only if needed)

**Checkpoint**: Recover test passes; Stories 1–2 still pass

---

## Phase 6: Polish & Cross-Cutting Concerns

- [x] T014 Update `AGENTS.md` Implemented / Remaining / Active for `007` shipped
- [x] T015 Run full suite: `go test ./processing/ ./gateway/ ./protocol/proto/event/v1/`

---

## Dependencies & Execution Order

- Phase 1 → Phase 2 → US1 → US2 → US3 → Polish
- US2 fixtures can be authored in parallel with US1 fixture (T006 ∥ T009)
- US3 depends on US1 happy path

## Implementation Strategy

1. Complete foundational deploy + StartEventID (T002–T005)
2. MVP: happy-path fixture + timer/message wins (US1)
3. Invalid deploy coverage (US2)
4. Recover (US3)
5. AGENTS.md + full suite
