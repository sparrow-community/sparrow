---
description: "Task list for incident (blocked execution state) increment"
---

# Tasks: Incident (blocked execution state)

**Input**: Design documents from `/specs/003-incident/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Included — constitution requires Event-intent tests for semantic changes; spec success criteria SC-001–SC-004.

**Organization**: Tasks are grouped by user story so each story can be implemented and tested independently.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1–US3
- Paths are repository-relative

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Wire-compatible protocol fields used by every story

- [x] T001 Add `INTENT_INCIDENT_OPENED` (14) and `INTENT_INCIDENT_RESOLVED` (15) to `Element.Intent` in `protocol/proto/event/v1/event.proto`
- [x] T002 [P] Add `job_fail_count` to `ActivityPayload` in `protocol/proto/event/v1/payloads.proto`
- [x] T003 [P] Add `incident_error_message` and `job_fail_count` to `Token`, plus `ResolveIncidentRequest`/`ResolveIncidentResponse` and `ResolveIncident` RPC in `protocol/proto/engine/v1/engine.proto`
- [x] T004 [P] Add `no_retry` to `FailJobRequest` in `protocol/proto/job/v1/job.proto`
- [x] T005 Regenerate Go with `cd protocol/proto && ./build.sh` (do not hand-edit `protocol/gen/go`)

**Checkpoint**: `buf lint` clean; existing `go test ./protocol/proto/event/v1/` still passes

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Deploy threshold compile, projection fields, and replay hooks every story needs; no user-visible incident behavior until US1

- [x] T006 Add `IncidentThreshold` to service-task deploy spec and parse optional `sparrow:failedJobIncidentThreshold` from BPMN extension elements in `processing/deploy/` (new `incident.go` or extend existing compile path in `processing/deploy/deploy.go`); default `3` when absent
- [x] T007 Extend `projection.Token` with `TokenBlocked` status constant, `JobFailCount`, and `IncidentErrorMessage` in `processing/projection/instance.go`
- [x] T008 Handle `INTENT_FAILED` fail-count increment, `INTENT_INCIDENT_OPENED`, `INTENT_INCIDENT_RESOLVED`, `INTENT_ACTIVATED` reset, and `INTENT_TERMINATED` clear in `ApplyEvent` in `processing/projection/instance.go`
- [x] T009 Add `defaultIncidentThreshold` and `IncidentThreshold(elementID)` lookup on `deploy.Deployment` in `processing/deploy/incident.go`

**Checkpoint**: Existing `go test ./processing/` passes (no incident fixtures yet)

---

## Phase 3: User Story 1 - Job failure opens an incident (Priority: P1) 🎯 MVP

**Goal**: Retriable fails accumulate until threshold (default 3) or `no_retry` opens an incident; blocked tokens visible in `GetInstance`; Complete/Activate/Fail on blocked token rejected

**Independent Test**: Quickstart Story 1 — fail Service Task job three times → `status=blocked`; sub-threshold fail still retriable; Complete on blocked token → REJECTION `INCIDENT_OPEN`

### Tests for User Story 1

- [x] T010 [P] [US1] Add `processing/testdata/m7_incident_service.bpmn` (start → Service Task with job type → end)
- [x] T011 [US1] Add `processing/incident_test.go` with sub-threshold retry, threshold incident open, `no_retry` immediate incident, Complete-while-blocked rejection, and Activate-skips-blocked cases

### Implementation for User Story 1

- [x] T012 [US1] Extend `Engine.Fail` in `processing/jobs.go` to accept `noRetry`, increment fail count, emit `INTENT_INCIDENT_OPENED` after `INTENT_FAILED` when count ≥ threshold or `noRetry`
- [x] T013 [US1] Reject `Fail` on blocked tokens with code `INCIDENT_OPEN` in `processing/jobs.go`
- [x] T014 [US1] Exclude blocked tokens from `Activate` job queue in `processing/jobs.go`
- [x] T015 [US1] Reject `Complete` on blocked tokens with code `INCIDENT_OPEN` in `processing/engine.go`
- [x] T016 [US1] Map `blocked` status, `incident_error_message`, and `job_fail_count` on tokens through `gateway/engine_server.go`

**Checkpoint**: Story 1 quickstart passes; `TestJobServiceFailThenActivate` still passes for sub-threshold fails

---

## Phase 4: User Story 2 - Resolve incident and retry work (Priority: P1)

**Goal**: Explicit `ResolveIncident` COMMAND closes incident and restores waiting job semantics; worker can activate and complete afterward

**Independent Test**: Quickstart Story 2 — resolve blocked token → waiting → activate → complete → process ends; resolve on non-blocked token → REJECTION `NO_INCIDENT`

### Tests for User Story 2

- [x] T017 [P] [US2] Add resolve-and-complete and `NO_INCIDENT` rejection tests to `processing/incident_test.go`
- [x] T018 [US2] Add `ResolveIncident` RPC test in `gateway/engine_server_test.go`

### Implementation for User Story 2

- [x] T019 [US2] Implement `Engine.ResolveIncident` with COMMAND/EVENT `INTENT_INCIDENT_RESOLVED` chain in `processing/engine.go`
- [x] T020 [US2] Wire `ResolveIncident` RPC in `gateway/engine_server.go`
- [x] T021 [US2] Pass `FailJobRequest.no_retry` through `gateway/job_server.go` to `Engine.Fail`

**Checkpoint**: Stories 1 and 2 quickstart paths pass; audit shows FAILED + INCIDENT_OPENED + INCIDENT_RESOLVED lifecycle

---

## Phase 5: User Story 3 - Recover with an open incident (Priority: P2)

**Goal**: Recover rebuilds blocked state and fail counts from the log; redrive completes Fail command chains without duplicate INCIDENT_OPENED

**Independent Test**: Quickstart Story 3 — open incident, new engine + Recover → same blocked state; resolve and complete after recover matches non-recover outcome

### Tests for User Story 3

- [x] T022 [P] [US3] Add Recover-rebuilds-blocked-state test to `processing/recover_test.go`
- [x] T023 [US3] Add redrive-after-partial-Fail-command test (COMMAND appended, INCIDENT_OPENED missing) to `processing/recover_test.go`

### Implementation for User Story 3

- [x] T024 [US3] Extend `recover.go` redrive for `INTENT_FAILED` commands to emit `INTENT_INCIDENT_OPENED` when policy triggers, using existing `alreadySeen` idempotency
- [x] T025 [US3] Ensure full log replay in Recover rebuilds `JobFailCount`, blocked status, and incident fields consistently in `processing/recover.go` and `processing/projection/instance.go`

**Checkpoint**: Story 3 quickstart passes; no duplicate incident EVENTs on redrive

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Edge cases, regression, and roadmap updates

- [x] T026 [P] Add `processing/testdata/m7_incident_boundary.bpmn` and boundary-supersedes-incident test in `processing/incident_test.go`
- [x] T027 [P] Add resolve-on-completed-instance and User-Task-fail-no-incident rejection tests in `processing/incident_test.go`
- [x] T028 Run full regression per quickstart: `go test ./processing/ ./gateway/ ./protocol/proto/event/v1/`
- [x] T029 Update `AGENTS.md`: move Incident from Remaining to Implemented; clear or advance active increment pointer

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 (proto generated) — **blocks all user stories**
- **User Story 1 (Phase 3)**: Depends on Phase 2 — MVP
- **User Story 2 (Phase 4)**: Depends on Phase 2; practically builds on US1 blocked state but tests can open incident inline
- **User Story 3 (Phase 5)**: Depends on Phase 2; needs US1 Fail/incident open behavior for meaningful Recover tests
- **Polish (Phase 6)**: Depends on Phases 3–5

### User Story Dependencies

- **US1 (P1)**: After Foundational — no dependency on US2/US3
- **US2 (P1)**: After Foundational — uses incident from US1 Fail path; independently testable by opening incident in test setup
- **US3 (P2)**: After US1 (incident open path) — Recover tests need persisted INCIDENT_OPENED facts

### Parallel Opportunities

- Phase 1: T002, T003, T004 in parallel after T001 enum slot is reserved
- Phase 3: T010 fixture in parallel with early Phase 2 if proto done
- Phase 4: T017 in parallel with T019 once US1 Fail path exists
- Phase 5: T022 fixture/setup parallel with T024 implementation planning
- Phase 6: T026 and T027 in parallel

---

## Parallel Example: User Story 1

```bash
# After Phase 2 completes:
# Fixture + proto mapping in parallel:
Task T010: "Add m7_incident_service.bpmn"
Task T016: "Map blocked token fields in gateway/engine_server.go"

# Then sequential core:
Task T012 → T013 → T014 → T015 → T011 (tests validate full story)
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001–T005)
2. Complete Phase 2: Foundational (T006–T009)
3. Complete Phase 3: User Story 1 (T010–T016)
4. **STOP and VALIDATE**: Quickstart Story 1; `go test ./processing/`
5. Demo blocked instance visible in `GetInstance` without resolve path

### Incremental Delivery

1. Setup + Foundational → proto and projection ready
2. US1 → blocked incident on job failure (MVP)
3. US2 → operator resolve loop complete
4. US3 → Recover/completeness gate
5. Polish → boundary edge, AGENTS.md, full regression

### Suggested MVP Scope

**T001–T016** (16 tasks): Phases 1–3 only. Delivers SC-001 and partial SC-004 (Complete-while-blocked). Resolve and Recover deferred to US2/US3.

---

## Notes

- Do not hand-edit `protocol/gen/go/*.pb.go`; always run `build.sh`
- Rejection codes: `INCIDENT_OPEN`, `NO_INCIDENT`, plus existing `INVALID_STATE` / `NOT_FOUND`
- Blocked tokens keep `job_type` for post-resolve activation
- Serial instance lock preserves single incident open and consistent fail count under concurrent Fail commands
