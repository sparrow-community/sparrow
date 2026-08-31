---
description: "Task list for multi-instance activities increment"
---

# Tasks: Multi-Instance Activities

**Input**: Design documents from `/specs/002-multi-instance/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Included — constitution requires Event-intent tests for semantic changes.

**Organization**: Tasks are grouped by user story so each story can be implemented and tested independently.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1–US5
- Paths are repository-relative

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Wire-compatible protocol fields used by every story

- [x] T001 Add `loop_instance_index`, `loop_total_instances`, and `loop_completed_instances` to `ActivityPayload` in `protocol/proto/event/v1/payloads.proto`
- [x] T002 [P] Add `loop_instance_index` to `Token` in `protocol/proto/engine/v1/engine.proto`
- [x] T003 Regenerate Go with `cd protocol/proto && ./build.sh` (do not hand-edit `protocol/gen/go`)

**Checkpoint**: `buf lint` clean; existing `go test ./protocol/proto/event/v1/` still passes

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Deploy compile, projection, handler effects, and loop helpers every story needs; no user-visible MI behavior until US1

- [x] T004 Compile `MultiInstanceSpec` from `element.MultiInstanceLoopCharacteristics` in `processing/deploy/multi_instance.go`; wire into `processing/deploy/deploy.go`; reject `complexBehaviorDefinition`, none/one behavior event refs, and MI on unsupported element types
- [x] T005 Extend `projection.Token` with `LoopInstanceIndex` and `projection.Instance` with `MultiInstanceLoops` in `processing/projection/instance.go`
- [x] T006 Apply loop payload fields and rebuild `MultiInstanceLoops` counters in `ApplyEvent` in `processing/projection/instance.go`
- [x] T007 Extend `handlers.Effect` with multi-instance spawn, join, and cancel-straggler fields in `processing/handlers/handler.go`
- [x] T008 Add cardinality/collection resolution and loop-counter variable injection helpers in `processing/multi_instance.go`
- [x] T009 Expose `nrOfInstances`, `nrOfActiveInstances`, `nrOfCompletedInstances`, and `loopCounter` to `processing/expr` evaluation from `processing/multi_instance.go`

**Checkpoint**: Existing `go test ./processing/` passes (no MI fixtures yet)

---

## Phase 3: User Story 1 - Parallel multi-instance with fixed cardinality (Priority: P1) 🎯 MVP

**Goal**: Parallel MI with numeric `loopCardinality` on User Task and Service Task; all inner instances must complete; single outgoing flow

**Independent Test**: Quickstart Story 1 — three waiting tokens on parallel user task; three jobs on parallel service task; Recover mid-loop

### Tests for User Story 1

- [x] T010 [P] [US1] Add `processing/testdata/m6_mi_parallel_cardinality.bpmn` and `processing/testdata/m6_mi_parallel_service.bpmn`
- [x] T011 [US1] Add `processing/multi_instance_test.go` with parallel user-task join, parallel service-task jobs, zero-cardinality immediate complete, and Recover with 1 of 3 completed

### Implementation for User Story 1

- [x] T012 [US1] Branch `UserTaskHandler` `OnEnter`/`OnComplete` for parallel MI in `processing/handlers/user_task.go`
- [x] T013 [US1] Branch `ServiceTaskHandler` `OnEnter`/`OnComplete` for parallel MI in `processing/handlers/service_task.go`
- [x] T014 [US1] Implement parallel inner spawn, per-inner `Complete`, all-complete join, and single outgoing flow in `processing/executor.go`
- [x] T015 [US1] Map `loop_instance_index` on tokens through `gateway/engine_server.go` and cover in `gateway/engine_server_test.go`
- [x] T016 [US1] Treat zero or negative evaluated cardinality as zero instances with immediate host completion in `processing/multi_instance.go`

**Checkpoint**: Story 1 quickstart passes including Recover; non-MI user/service task tests still pass

---

## Phase 4: User Story 2 - Sequential multi-instance (Priority: P1)

**Goal**: Sequential MI with fixed cardinality; at most one active inner instance; ordered activation

**Independent Test**: Quickstart Story 2 — one waiting token at a time; three ordered completes finish the process

### Tests for User Story 2

- [x] T017 [P] [US2] Add `processing/testdata/m6_mi_sequential_cardinality.bpmn`
- [x] T018 [US2] Add sequential-only-waiting-token and ordered-completion tests to `processing/multi_instance_test.go`

### Implementation for User Story 2

- [x] T019 [US2] Activate the next `loop_instance_index` only after inner `COMPLETED` when `isSequential=true` in `processing/executor.go`
- [x] T020 [US2] Ensure completed inner indices do not re-open and sequential counters stay consistent in `processing/projection/instance.go`

**Checkpoint**: Story 1 still passes; Story 2 sequential fixture passes

---

## Phase 5: User Story 3 - Collection-driven input (Priority: P2)

**Goal**: `loopDataInputRef` drives instance count; `inputDataItem` per iteration; optional `loopDataOutputRef` assembly; empty collection completes immediately

**Independent Test**: Quickstart Story 3 — `items=["a","b"]` → two inner instances with element vars; empty collection fixture

### Tests for User Story 3

- [x] T021 [P] [US3] Add `processing/testdata/m6_mi_collection.bpmn` and `processing/testdata/m6_mi_collection_empty.bpmn`
- [x] T022 [US3] Add collection input, output collection, and empty-collection tests to `processing/multi_instance_test.go`

### Implementation for User Story 3

- [x] T023 [US3] Resolve collection size and bind `inputDataItem` variable per inner instance in `processing/multi_instance.go`
- [x] T024 [US3] Append `outputDataItem` values and write `loopDataOutputRef` on host complete in `processing/executor.go`
- [x] T025 [US3] Treat missing or non-array collection variables as empty (zero instances) in `processing/multi_instance.go`

**Checkpoint**: Stories 1–2 still pass; Story 3 collection fixtures pass

---

## Phase 6: User Story 4 - Completion condition before all instances finish (Priority: P2)

**Goal**: Custom `completionCondition` and `behavior=One` allow early join; remaining inner instances cancelled

**Independent Test**: Quickstart Story 4 — 2 of 5 parallel completes finish the activity; no stray waiting tokens

### Tests for User Story 4

- [x] T026 [P] [US4] Add `processing/testdata/m6_mi_completion_early.bpmn`
- [x] T027 [US4] Add early-completion, straggler-cancel, and default-all-complete regression tests to `processing/multi_instance_test.go`

### Implementation for User Story 4

- [x] T028 [US4] Evaluate `completionCondition` and map `behavior` One/All to default expressions after each inner complete in `processing/executor.go`
- [x] T029 [US4] Cancel remaining active inner tokens and complete host once when completion condition is satisfied in `processing/executor.go`

**Checkpoint**: Stories 1–3 still pass; Story 4 early-exit fixture passes

---

## Phase 7: User Story 5 - Multi-instance embedded sub-process (Priority: P3)

**Goal**: MI on embedded Sub-Process; parallel or sequential inner scopes; per-scope join then outer loop join

**Independent Test**: Quickstart Story 5 — two parallel sub-process scopes; both must complete before outer MI sub-process completes once

### Tests for User Story 5

- [x] T030 [P] [US5] Add `processing/testdata/m6_mi_subprocess.bpmn`
- [x] T031 [US5] Add MI sub-process parallel join and partial-completion-waits tests to `processing/multi_instance_test.go`

### Implementation for User Story 5

- [x] T032 [US5] Integrate MI loop host with `SubProcessHandler` `OnEnter`/`OnComplete` in `processing/handlers/sub_process.go`
- [x] T033 [US5] Wire inner scope `COMPLETED` to loop counters and outer host join in `processing/executor.go`

**Checkpoint**: Stories 1–4 still pass; embedded non-MI SubProcess tests still pass

---

## Phase 8: Polish & Cross-Cutting Concerns

- [x] T034 [P] Cancel all inner instances and loop host when an interrupting boundary fires on a multi-instance activity in `processing/handlers/boundary_event.go`
- [x] T035 [P] Document multi-instance semantics (parallel/sequential, collection, completion, recover) in `processing/README.md`
- [x] T036 [P] Update implemented-element snapshot and active increment pointer in `AGENTS.md`
- [x] T037 Run `go test ./processing/ ./gateway/ ./protocol/proto/event/v1/` and all paths in `specs/002-multi-instance/quickstart.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies
- **Foundational (Phase 2)**: Depends on T001–T003
- **US1 (Phase 3)**: Depends on Phase 2 — MVP
- **US2 (Phase 4)**: Depends on US1 (shared executor join path; sequential is a mode flag on the same machinery)
- **US3 (Phase 5)**: Depends on US1 (parallel spawn); collection replaces cardinality resolution only
- **US4 (Phase 6)**: Depends on US1 (parallel with multiple active inners); early cancel extends join
- **US5 (Phase 7)**: Depends on US1–US2 (loop counters + sequential mode); Sub-Process scope integration is additive
- **Polish**: After the stories you intend to ship

### User Story Dependencies

- **US1 (P1)**: After Phase 2 — core parallel MI
- **US2 (P1)**: After US1 — sequential activation
- **US3 (P2)**: After US1 — collection input/output
- **US4 (P2)**: After US1 — early completion (needs multiple parallel inners)
- **US5 (P3)**: After US1–US2 — MI embedded Sub-Process

### Parallel Opportunities

- T001 and T002
- T010 fixtures vs T012 handler work (after Phase 2)
- T017 vs T019 (fixture vs executor, after US1)
- T021 vs T023 (fixture vs collection resolver, after US1)
- T026 vs T028 (fixture vs completion eval, after US1)
- T030 vs T032 (fixture vs sub_process handler, after US1)
- T034 / T035 / T036 documentation and boundary polish

---

## Parallel Example: User Story 1

```bash
# Fixtures and gateway mapping can proceed beside handler work once Phase 2 exists:
Task: "Add processing/testdata/m6_mi_parallel_cardinality.bpmn"
Task: "Map loop_instance_index in gateway/engine_server.go"
Task: "Branch UserTaskHandler for parallel MI in processing/handlers/user_task.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Phase 1 + Phase 2
2. Phase 3 (US1)
3. Stop and validate quickstart Story 1 including Recover

### Incremental Delivery

1. US1 → parallel MI (user + service task)
2. US2 → sequential mode
3. US3 → collection input/output
4. US4 → completion condition / early exit
5. US5 → MI embedded Sub-Process

### Notes

- Do not implement deferred items listed in `AGENTS.md` and `spec.md` (MI Call Activity, complex behavior, cluster, product suite, …)
- Prefer extending handlers and `processing/multi_instance.go` over new Engine lifecycle branches
- `Complete` stays token-id based; inner MI work completes via distinct inner token ids
- Commit after each story checkpoint if the user asks for commits
