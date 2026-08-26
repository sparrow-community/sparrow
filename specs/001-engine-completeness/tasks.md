---
description: "Task list for engine completeness next increment"
---

# Tasks: Engine Completeness (next increment)

**Input**: Design documents from `/specs/001-engine-completeness/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Included — constitution requires Event-intent tests for semantic changes.

**Organization**: Tasks are grouped by user story so each story can be implemented and tested independently.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1–US4
- Paths are repository-relative

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Wire-compatible protocol fields used by later stories

- [ ] T001 Add optional parent/call fields to `ProcessPayload` and `called_process_instance_id` to `ActivityPayload` in `protocol/proto/event/v1/payloads.proto`
- [ ] T002 [P] Add `process_id` / `process_version` selection on `CreateInstanceRequest`, parent/called fields on `Instance`/`Token`, and deploy response ids in `protocol/proto/engine/v1/engine.proto`
- [ ] T003 Regenerate Go with `cd protocol/proto && ./build.sh` (do not hand-edit `protocol/gen/go`)

**Checkpoint**: `buf lint` clean; existing `go test ./protocol/proto/event/v1/` still passes

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Projection and deploy hooks every story needs; no user-visible Call Activity behavior change yet except dropping v1 uniqueness so US1 can deploy

- [ ] T004 Extend `projection.Instance` and `projection.Token` with parent/called ids in `processing/projection/instance.go` and apply them from the new payload fields in `ApplyEvent`
- [ ] T005 Compile Call Activity IO mapping structs (may be empty) and drop “one Call Activity per called process” / “no nested Call Activity” checks in `processing/deploy/call_activity.go` and `processing/deploy/deploy.go`
- [ ] T006 Add a post-lock child-instance publication kind (create child / resume parent) next to existing message/signal publication in `processing/executor.go` and `processing/handlers/handler.go` without changing inline SubProcess behavior

**Checkpoint**: Existing `go test ./processing/` passes (inline Call Activity tests may still expect same-instance tokens until US1)

---

## Phase 3: User Story 1 - Called process as its own instance (Priority: P1) 🎯 MVP

**Goal**: Call Activity starts a distinct child instance; caller waits; child complete resumes caller; Recover keeps the link

**Independent Test**: Quickstart Story 1 fixtures; two Call Activities naming the same child deploy and run

### Tests for User Story 1

- [ ] T007 [P] [US1] Add `processing/testdata/m5_call_instance.bpmn` (caller + child user task) and `processing/testdata/m5_call_instance_twice.bpmn` (two Call Activities → same child)
- [ ] T008 [US1] Rewrite `processing/call_activity_test.go` for child `process_instance_id`, parent host token, complete-child-then-caller, dual-call, interrupting boundary, and Recover mid-call

### Implementation for User Story 1

- [ ] T009 [US1] Change `handlers.CallActivityHandler` in `processing/handlers/call_activity.go` to park the host token and request child start (payload: `called_process_instance_id`) instead of `EnterChild` on the same instance
- [ ] T010 [US1] Implement child CreateInstance under the child lock and parent Complete/Terminate of the Call Activity after child PROCESS COMPLETED/TERMINATED in `processing/engine.go` / `processing/executor.go`
- [ ] T011 [US1] Map parent/called fields through `gateway/engine_server.go` `GetInstance` snapshots
- [ ] T012 [US1] Update Call Activity semantics in `processing/DESIGN.md` §5.1 / §8 to match independent instances (remove v1 inline debt note)

**Checkpoint**: Story 1 quickstart passes; embedded SubProcess tests still pass

---

## Phase 4: User Story 2 - Map variables into and out of the called instance (Priority: P1)

**Goal**: Name-to-name IO associations at start and complete; no mapping ⇒ no copy

**Independent Test**: Quickstart Story 2

### Tests for User Story 2

- [ ] T013 [P] [US2] Add `processing/testdata/m5_call_io.bpmn` and `processing/testdata/m5_call_no_io.bpmn`
- [ ] T014 [US2] Add mapping tests in `processing/call_activity_test.go` (input only, output only, missing source skipped, no mapping)

### Implementation for User Story 2

- [ ] T015 [US2] Resolve `DataInputAssociations` / `DataOutputAssociations` into compiled mappings in `processing/deploy/call_activity.go`; reject unknown structure at deploy
- [ ] T016 [US2] Apply inputs when starting the child and outputs when completing the Call Activity in `processing/handlers/call_activity.go` / `processing/executor.go`

**Checkpoint**: Story 1 still passes with unmapped fixtures; Story 2 mappings hold

---

## Phase 5: User Story 3 - Error Event Sub-Process (Priority: P2)

**Goal**: Error starts on Event Sub-Process (process and embedded sub-process); interrupting and non-interrupting; code vs catch-all

**Independent Test**: Quickstart Story 3

### Tests for User Story 3

- [ ] T017 [P] [US3] Add `processing/testdata/m5_error_esp.bpmn`, `m5_error_esp_ni.bpmn`, `m5_error_esp_code_miss.bpmn`, `m5_error_esp_nested.bpmn`
- [ ] T018 [US3] Add `processing/error_event_subprocess_test.go` covering the four fixtures and Recover while the ESP is armed

### Implementation for User Story 3

- [ ] T019 [P] [US3] Extend `deploy.CatchKind` / `EventSubProcess` with error in `processing/deploy/timer.go` and `processing/deploy/event_subprocess.go`; set `EventPayload.error_code` on START_EVENT ACTIVATED
- [ ] T020 [US3] Route unmatched/bubbling errors to matching ESP in `processing/errors.go` and `processing/event_subprocess.go` (innermost: activity boundary, then scope ESP)
- [ ] T021 [US3] Arm error ESP when opening a process or embedded sub-process in `processing/event_subprocess.go` (same place message/timer/signal arms are created)

**Checkpoint**: Existing error boundary / ThrowError tests still pass; Story 3 fixtures pass

---

## Phase 6: User Story 4 - Keep process versions side by side (Priority: P3)

**Goal**: Same process id, incrementing revisions; start latest or explicit version; running instances stay bound

**Independent Test**: Quickstart Story 4

### Tests for User Story 4

- [ ] T022 [P] [US4] Add `processing/testdata/m5_version_v1.bpmn` and `processing/testdata/m5_version_v2.bpmn`
- [ ] T023 [US4] Add `processing/version_test.go` (latest start, explicit old version, unknown version rejected, in-flight instance stays on v1)

### Implementation for User Story 4

- [ ] T024 [US4] Assign `process_version` per process id and keep a revision index in `processing/deploy/deploy.go`; load it on Recover from `deploy.Store`
- [ ] T025 [US4] Extend `Engine.CreateInstance` in `processing/engine.go` to accept `process_id` + optional version; keep `deployment_id` path
- [ ] T026 [US4] Map new Deploy/CreateInstance fields in `gateway/engine_server.go` and cover them in `gateway/engine_server_test.go`

**Checkpoint**: Stories 1–3 still pass when started via `deployment_id`

---

## Phase 7: Polish & Cross-Cutting Concerns

- [x] T027 [P] Merge roadmap / implemented snapshot into `AGENTS.md` and point next steps at this feature’s `tasks.md` (no separate `STATUS.md`)
- [x] T028 [P] Keep Spec Kit constitution / active spec links in `AGENTS.md` (do not duplicate DESIGN)
- [ ] T029 Run `go test ./processing/ ./gateway/ ./protocol/proto/event/v1/` and the Story 1 Recover path from `specs/001-engine-completeness/quickstart.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **Foundational (Phase 2)**: Depends on T001–T003
- **US1 (Phase 3)**: Depends on Phase 2 — MVP
- **US2 (Phase 4)**: Depends on US1 (mappings apply at child start/complete)
- **US3 (Phase 5)**: Depends on Phase 2 only; can proceed in parallel with US1/US2 if staffing allows
- **US4 (Phase 6)**: Depends on Phase 2 only; CreateInstance signature should land after US1 if US1 already calls CreateInstance internally
- **Polish**: After the stories you intend to ship

### User Story Dependencies

- **US1 (P1)**: After Phase 2
- **US2 (P1)**: After US1
- **US3 (P2)**: After Phase 2; independent of Call Activity stories
- **US4 (P3)**: After Phase 2; coordinate with US1’s internal CreateInstance

### Parallel Opportunities

- T001 and T002
- T007 fixtures vs T009 handler (after T005)
- T013 fixtures vs T015 compile
- T017 / T019 vs Call Activity work
- T022 vs T024
- T027 / T028 documentation

---

## Parallel Example: User Story 1

```bash
# Fixtures and DESIGN can proceed beside handler work once T005/T006 exist:
Task: "Add processing/testdata/m5_call_instance.bpmn"
Task: "Park host token in processing/handlers/call_activity.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Phase 1 + Phase 2
2. Phase 3 (US1)
3. Stop and validate quickstart Story 1 including Recover

### Incremental Delivery

1. US1 → independent called instance
2. US2 → mappings
3. US3 → error Event Sub-Process
4. US4 → revision coexistence

### Notes

- Do not implement deferred items listed in `AGENTS.md` (cluster, product suite, cross-deployment Call Activity, live migration, multi-instance, …)
- Prefer extending handlers over adding Engine lifecycle branches
- Commit after each story checkpoint if the user asks for commits
