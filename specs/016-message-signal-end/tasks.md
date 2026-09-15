---
description: "Task list for Message End and Signal End"
---

# Tasks: Message End and Signal End

**Input**: Design documents from `specs/016-message-signal-end/`

## Phase 1: Setup

- [x] T001 Point `.specify/feature.json` at `specs/016-message-signal-end`

## Phase 2: Foundational

- [x] T002 Index message/signal ends; reject mixed defs in `processing/deploy/end.go` and `processing/deploy/deploy.go`

## Phase 3: User Story 1 - Message end (P1)

- [x] T003 [US1] Publish on message end enter in `processing/handlers/end_event.go`
- [x] T004 [P] [US1] Fixtures `m20_message_end.bpmn` / `m20_message_end_to_catch.bpmn` + tests in `processing/message_signal_end_test.go`

## Phase 4: User Story 2 - Signal end (P1)

- [x] T005 [US2] Publish on signal end enter in `processing/handlers/end_event.go`
- [x] T006 [P] [US2] Fixtures `m20_signal_end.bpmn` / `m20_signal_end_to_catch.bpmn` + tests

## Phase 5: User Story 3 - SubProcess ends (P2)

- [x] T007 [P] [US3] Fixtures `m20_subprocess_message_end.bpmn` / `m20_subprocess_signal_end.bpmn` + tests

## Phase 6: User Story 4 - Recover (P2)

- [x] T008 [US4] Recover tests for message/signal end delivery in `processing/message_signal_end_test.go`

## Phase 7: Polish

- [x] T009 Update `AGENTS.md` and `processing/README.md`; run `go test ./processing/`
