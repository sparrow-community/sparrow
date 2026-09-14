---
description: "Task list for Process-Level Typed Start"
---

# Tasks: Process-Level Typed Start

**Input**: Design documents from `specs/015-process-typed-start/`

**Tests**: Included — spec SC/quickstart require fixtures and Recover coverage.

## Phase 1: Setup

- [x] T001 Point `.specify/feature.json` at `specs/015-process-typed-start`

## Phase 2: Foundational

- [x] T002 Add process-level start classification, indexes, and `NoneStartEventID` in `processing/deploy/start.go`
- [x] T003 Index/validate process-level typed starts; reject process-level error/escalation/unsupported defs in `processing/deploy/deploy.go`
- [x] T004 Require Call Activity called process none or instantiate entry in `processing/deploy/call_activity.go`
- [x] T005 Extract shared instance mint helper `createInstanceAt` and route CreateInstance to none/instantiate only in `processing/engine.go`

**Checkpoint**: Foundation ready — typed indexes exist; CreateInstance rejects typed-only

---

## Phase 3: User Story 1 - Message start (P1) 🎯 MVP

**Goal**: Unmatched PublishMessage creates instance at message start; CreateInstance rejected for message-only.

**Independent Test**: `m19_message_start.bpmn` — CreateInstance fails; PublishMessage → User Task wait → Complete.

- [x] T006 [US1] Wire message-start creation after unmatched PublishMessage in `processing/messages.go`
- [x] T007 [P] [US1] Fixture `processing/testdata/m19_message_start.bpmn` + tests in `processing/typed_start_test.go`

**Checkpoint**: Message start MVP works

---

## Phase 4: User Story 2 - Timer start (P1)

**Goal**: Deploy arms timer start; FireDue mints instance when due.

**Independent Test**: `m19_timer_start.bpmn` — CreateInstance fails; FireDue creates and progresses.

- [x] T008 [US2] Arm/persist process-level timer-start schedules; FireDue creates instances in `processing/timers.go` (+ runtime helper as needed)
- [x] T009 [P] [US2] Fixture `processing/testdata/m19_timer_start.bpmn` + FireDue tests in `processing/typed_start_test.go`

---

## Phase 5: User Story 3 - Signal start (P1)

**Goal**: Unmatched PublishSignal creates instance at signal start.

**Independent Test**: `m19_signal_start.bpmn` — PublishSignal creates and progresses.

- [x] T010 [US3] Wire signal-start creation after unmatched PublishSignal in `processing/signals.go`
- [x] T011 [P] [US3] Fixture `processing/testdata/m19_signal_start.bpmn` + tests in `processing/typed_start_test.go`

---

## Phase 6: User Story 4 - Conditional start (P1)

**Goal**: EvaluateConditionalStarts creates instance when condition true.

**Independent Test**: false → no instance; true → progresses.

- [x] T012 [US4] Add `EvaluateConditionalStarts` in `processing/conditional_start.go`
- [x] T013 [P] [US4] Fixture `processing/testdata/m19_conditional_start.bpmn` + tests in `processing/typed_start_test.go`

---

## Phase 7: User Story 5 - None + mixed alternatives (P1)

**Goal**: None CreateInstance unchanged; none+message alternatives work.

**Independent Test**: `m19_none_and_message_start.bpmn` — CreateInstance via none; PublishMessage via message start.

- [x] T014 [US5] Ensure CreateInstance uses none when present among alternatives in `processing/engine.go` / `processing/deploy/start.go`
- [x] T015 [P] [US5] Fixture `processing/testdata/m19_none_and_message_start.bpmn` + tests; confirm instantiate-EBG still needs CreateInstance

---

## Phase 8: User Story 6 - Process-level error start reject (P2)

**Goal**: Process-level error start Deploy rejected; ESP error suite green.

**Independent Test**: `m19_error_start_reject.bpmn` Deploy fails; existing ESP error tests pass.

- [x] T016 [P] [US6] Fixture `processing/testdata/m19_error_start_reject.bpmn` + Deploy reject test in `processing/typed_start_test.go`

---

## Phase 9: User Story 7 - Recover (P2)

**Goal**: Recover preserves timer-start arms; message/timer/signal paths match continuous.

**Independent Test**: Recover mid-path then finish matches continuous for message/timer/signal.

- [x] T017 [US7] Rebuild timer-start schedules on Open/Recover in `processing/recover.go` / runtime path
- [x] T018 [US7] Recover tests for message/timer/signal starts in `processing/typed_start_test.go`

---

## Phase 10: Polish

- [x] T019 Update `AGENTS.md` Supported/Planned and `processing/README.md` for typed starts
- [x] T020 Run `go test ./processing/` and fix regressions

---

## Dependencies & Execution Order

- Setup → Foundational (blocks all stories)
- US1 (message) first MVP; US2–US5 parallelizable after foundation (shared test file → sequential tests preferred)
- US6 can run anytime after T003
- US7 after US1–US3 mint paths exist
- Polish last

### Parallel opportunities

- Fixture files T007/T009/T011/T013/T015/T016 are [P] across files
- Implementation touching `messages.go` / `signals.go` / `timers.go` / `conditional_start.go` can proceed in parallel after T005

### MVP

T001–T007 (foundation + message start)

## Implementation Strategy

1. Foundation: classify starts, reject invalid, mint helper, CreateInstance gate
2. Message start MVP → validate
3. Timer / signal / conditional / mixed
4. Error reject + Recover
5. Docs + full suite
