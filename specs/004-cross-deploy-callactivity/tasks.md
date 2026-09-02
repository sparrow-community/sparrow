---
description: "Task list for cross-deployment Call Activity increment"
---

# Tasks: Cross-deployment Call Activity

**Input**: Design documents from `/specs/004-cross-deploy-callactivity/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Included — constitution requires Event-intent tests for semantic changes; spec success criteria SC-001–SC-005.

**Organization**: Tasks are grouped by user story so each story can be implemented and tested independently.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1–US3
- Paths are repository-relative

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Confirm baseline before semantic changes

- [x] T001 Run SC-005 baseline: `go test ./processing/ -run CallActivity -count=1` and record passing state for `m4_call_activity.bpmn` / `m5_call_*.bpmn` in `processing/call_activity_test.go`

**Checkpoint**: Existing same-deployment Call Activity tests pass

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Dual-mode deploy compile, publication fields, callee resolver hook, and child start on callee deployment — required by every user story

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T002 Split `validateCallActivity` into ref extraction + embedded catalog validation; add `ExternalCallee bool` and allow empty `StartEventID` for external refs in `processing/deploy/call_activity.go`
- [x] T003 Update `indexCallActivities` in `processing/deploy/deploy.go` to skip `calledProcesses` embed and callee scope indexing when `ExternalCallee` is true
- [x] T004 Add `CalledDeploymentID` to `handlers.Publication` in `processing/handlers/handler.go`
- [x] T005 Add `SetCalleeResolver(func(processID string) (deploymentID string, err error))` in `processing/handlers/call_activity.go`
- [x] T006 Wire callee resolver from `processing` init (e.g. `call_child.go` `init`) to `Engine.resolveDeploymentLocked("", processID, 0)` in `processing/engine.go`
- [x] T007 Update `startCalledInstance` in `processing/call_child.go` to load callee via `p.CalledDeploymentID`, use callee `deployment_id`/`version` on child COMMAND/projection, resolve `StartEventID` from compile or `calleeDep.StartEventID()`, and call `Enter`/`emitEventSubProcessStartArms` with `calleeDep`

**Checkpoint**: Caller-only BPMN deploys without embedded callee; embedded same-file deploy unchanged; `go test ./processing/` passes

---

## Phase 3: User Story 1 - Call a separately deployed process (Priority: P1) 🎯 MVP

**Goal**: Two deployments; Call Activity spawns child with callee deployment id; caller waits; callee-missing → REJECTION `NOT_FOUND` without waiting projection

**Independent Test**: Quickstart Story 1 — deploy callee then caller; `GetInstance` shows distinct deployment ids and linkage; caller-only deploy reaches Call Activity → REJECTION

### Tests for User Story 1

- [x] T008 [P] [US1] Add `processing/testdata/m8_cross_call_callee.bpmn` (simple callee process with user task)
- [x] T009 [P] [US1] Add `processing/testdata/m8_cross_call_caller.bpmn` (caller-only with Call Activity referencing callee process id)
- [x] T010 [US1] Add cross-deploy spawn, child completion resumes caller, caller-only deploy success, and callee-not-deployed `NOT_FOUND` tests in `processing/call_activity_test.go`

### Implementation for User Story 1

- [x] T011 [US1] Pre-resolve external callee in `CallActivityHandler.OnEnter` before emitting ACTIVATING/ACTIVATED; return `NOT_FOUND` with no records when resolve fails in `processing/handlers/call_activity.go`
- [x] T012 [US1] Set `Publication.CalledDeploymentID` on enter (external: resolved callee id; embedded: caller deployment id) in `processing/handlers/call_activity.go`
- [x] T013 [US1] Ensure enter failure on external resolve surfaces COMMAND `REJECTION` without caller token `waiting` on Call Activity in `processing/engine.go` (adjust only if existing path does not reject cleanly)

**Checkpoint**: Story 1 quickstart passes; SC-001 and SC-004 satisfied

---

## Phase 4: User Story 2 - Variable mapping across deployment boundary (Priority: P1)

**Goal**: Input/output mappings on Call Activity work identically when caller and callee are in different deployments

**Independent Test**: Quickstart Story 2 — mapped IO round-trip; unmapped caller vars absent on child; caller vars unchanged when no mappings

### Tests for User Story 2

- [x] T014 [P] [US2] Add `processing/testdata/m8_cross_call_io_callee.bpmn` and `processing/testdata/m8_cross_call_io_caller.bpmn` with input/output associations
- [x] T015 [US2] Add cross-deploy IO mapping assertions (input copy, output copy, no-mapping case) to `processing/call_activity_test.go`

### Implementation for User Story 2

- [x] T016 [US2] Verify `applyMappings` input path in `startCalledInstance` reads caller `CallActivitySpec` from parent deployment in `processing/call_child.go` (fix if parent dep conflated with callee)
- [x] T017 [US2] Verify `resumeParentCall` output mapping reads child variables and applies caller `CallActivitySpec.Outputs` in `processing/call_child.go` (fix if callee dep used for mapping lookup)

**Checkpoint**: Story 2 quickstart passes; SC-002 satisfied

---

## Phase 5: User Story 3 - Recover with cross-deployment parent and child (Priority: P2)

**Goal**: Recover rebuilds waiting parent and active child with correct distinct deployment ids; redrive does not duplicate child or bind newer callee revision

**Independent Test**: Quickstart Story 3 — Recover after mid-call crash; complete child; parent finishes; redrive after partial child create

### Tests for User Story 3

- [x] T018 [P] [US3] Add Recover-rebuilds-cross-deploy-call test to `processing/recover_test.go`
- [x] T019 [US3] Add redrive-after-partial-child-create test (parent ACTIVATED, child COMMAND partial) to `processing/recover_test.go`

### Implementation for User Story 3

- [x] T020 [US3] Ensure `startCalledInstance` redrive uses `Publication.CalledDeploymentID` and preallocated `ChildInstanceID` without re-resolving latest callee revision in `processing/call_child.go` and `processing/recover.go`
- [x] T021 [US3] Verify full log replay projects parent and child `deployment_id` fields correctly when they differ in `processing/projection/instance.go` (adjust only if replay gap found)

**Checkpoint**: Story 3 quickstart passes; SC-003 satisfied

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Regression, roadmap, and quickstart validation

- [x] T022 [P] Run SC-005 regression: `go test ./processing/ -run CallActivity -count=1` — all `m4`/`m5` fixtures unchanged
- [x] T023 Run full suite per quickstart: `go test ./processing/ ./gateway/ ./protocol/proto/event/v1/`
- [x] T024 Update `AGENTS.md`: move Cross-deployment CallActivity from Remaining to Implemented; clear or advance active increment pointer

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Setup — **BLOCKS all user stories**
- **User Story 1 (Phase 3)**: Depends on Foundational — MVP
- **User Story 2 (Phase 4)**: Depends on Foundational; independently testable after US1 runtime path exists (IO uses same child start)
- **User Story 3 (Phase 5)**: Depends on Foundational + US1 cross-deploy call producing two deployment ids in log
- **Polish (Phase 6)**: Depends on US1–US3 (or at minimum US1 for partial ship)

### User Story Dependencies

```text
Phase 2 (Foundational)
    ↓
Phase 3 (US1) ──→ Phase 4 (US2)     [US2 tests IO on cross-deploy path]
    ↓
Phase 5 (US3)                         [Recover needs cross-deploy ledger]
    ↓
Phase 6 (Polish)
```

- **US1**: No dependency on US2/US3
- **US2**: Logical dependency on US1 child-start path; fixtures are independent files
- **US3**: Needs US1 behavior in log; no US2 requirement

### Within Each User Story

- Fixtures before integration tests
- Handler/runtime changes before e2e tests that assert them
- Implementation checkpoint before next story

### Parallel Opportunities

- **Phase 2**: T004 (`handler.go`) ∥ T002–T003 (`deploy/`) after T002 starts
- **US1**: T008 ∥ T009 (fixtures); T011–T013 sequential on `call_activity.go`
- **US2**: T014 (fixtures) while US1 tests run
- **US3**: T018 ∥ T019 (different test functions in `recover_test.go`)
- **Polish**: T022 ∥ T024

---

## Parallel Example: User Story 1

```bash
# Fixtures in parallel:
Task T008: "Add processing/testdata/m8_cross_call_callee.bpmn"
Task T009: "Add processing/testdata/m8_cross_call_caller.bpmn"

# After T011–T012 land, run:
go test ./processing/ -run 'CrossDeploy|CallActivity' -count=1
```

---

## Parallel Example: User Story 3

```bash
# Recover tests in parallel (same file, different functions — coordinate merge):
Task T018: "Recover-rebuilds-cross-deploy-call test in processing/recover_test.go"
Task T019: "Redrive partial child create test in processing/recover_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001)
2. Complete Phase 2: Foundational (T002–T007) — **critical**
3. Complete Phase 3: User Story 1 (T008–T013)
4. **STOP and VALIDATE**: Quickstart Story 1 + SC-004 negative case
5. Ship/demo cross-deploy Call Activity without IO/Recover stories if time-boxed

### Incremental Delivery

1. Setup + Foundational → compile + child-start infrastructure ready
2. US1 → independently testable cross-deploy spawn (MVP)
3. US2 → IO mapping parity across deployments
4. US3 → Recover + redrive safety
5. Polish → SC-005 regression + AGENTS.md

### Parallel Team Strategy

1. One developer: Foundational → US1 → US2 → US3 → Polish (sequential)
2. Two developers after Phase 2:
   - Dev A: US1 implementation (T011–T013)
   - Dev B: US1 fixtures + tests (T008–T010), then US2 fixtures (T014)
3. US3 after US1 merges

---

## Notes

- No new `EngineService` RPC; optional `ActivityPayload.called_process_id` proto field is **out of scope** unless audit gap found during implement
- Embedded same-file Call Activity must remain byte-for-byte behavior equivalent (SC-005)
- Both caller and callee deployment artifacts required in `deploy.Store` for Recover (documented in quickstart)
- `[P]` tasks touch different files; same-file test tasks should merge sequentially
