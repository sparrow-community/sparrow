---
description: "Task list for Terminate End and Remaining Tasks"
---

# Tasks: Terminate End and Remaining Tasks

**Input**: Design documents from `specs/014-terminate-and-tasks/`

## Phase 1: Setup

- [x] T001 Point `.specify/feature.json` at `specs/014-terminate-and-tasks`

## Phase 2: Foundational

- [x] T002 Index Manual/Receive/Send/Business Rule + terminate ends; reject abstract task and instantiate receive in `processing/deploy/deploy.go`
- [x] T003 Add `TerminateScope` on `processing/handlers/handler.go` Effect; register new handlers

## Phase 3: User Story 1 - Terminate end (P1)

- [x] T004 [US1] Terminate end OnEnter + executor scope cancel in `processing/handlers/end_event.go` and `processing/executor.go`
- [x] T005 [US1] Fixtures `processing/testdata/m18_terminate_parallel.bpmn` and `m18_terminate_subprocess.bpmn` plus tests in `processing/terminate_task_test.go`

## Phase 4: User Story 2 - Manual Task (P1)

- [x] T006 [US2] Manual Task handler in `processing/handlers/manual_task.go` and projection wait
- [x] T007 [US2] Fixture `m18_manual_task.bpmn` + Complete/Recover tests

## Phase 5: User Story 3 - Receive Task (P1)

- [x] T008 [US3] Receive Task handler + PublishMessage waiter fix in `processing/handlers/receive_task.go` and `processing/messages.go`
- [x] T009 [US3] Fixture `m18_receive_task.bpmn` + PublishMessage/Recover/instantiate-reject tests

## Phase 6: User Story 4 - Send Task (P1)

- [x] T010 [US4] Send Task handler in `processing/handlers/send_task.go`
- [x] T011 [US4] Fixtures `m18_send_task.bpmn` / `m18_send_to_receive.bpmn` + tests

## Phase 7: User Story 5 - Business Rule Task (P1)

- [x] T012 [US5] Business Rule handler + job notify/ThrowError/incident in `processing/handlers/business_rule_task.go`, `processing/engine.go`, `processing/errors.go`, `processing/recover.go`
- [x] T013 [US5] Fixture `m18_business_rule_task.bpmn` + Activate/Complete/Recover tests

## Phase 8: Polish

- [x] T014 Update `AGENTS.md` and `processing/README.md`; run `go test ./processing/`
