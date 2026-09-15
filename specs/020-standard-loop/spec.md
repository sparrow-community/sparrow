# Feature Specification: Standard Loop Characteristics

**Feature Branch**: `020-standard-loop`

**Created**: 2026-09-15

**Status**: Draft

**Input**: User description: "继续 Planned P2：standardLoopCharacteristics。Deploy 接受标准循环；testBefore / loopCondition / loopMaximum；活动重复执行直到条件结束；与 multiInstance 互斥；等待类任务与 Recover。"

## User Scenarios & Testing *(mandatory)*

Authors attach `standardLoopCharacteristics` so an activity repeats. Deploy currently rejects it. This increment executes standard loops with BPMN test-before / test-after semantics and an optional maximum.

### User Story 1 - Do-while style loop (Priority: P1)

An author models a User Task (or Manual/abstract Task) with `testBefore="false"`, a `loopCondition` on a process variable, and optional `loopMaximum`. The instance waits each iteration; Complete updates variables; while the condition holds (and under max), the same activity waits again; then the token leaves.

**Why this priority**: Default BPMN `testBefore` is false; covers the common “repeat until done” pattern.

**Independent Test**: Fixture with condition `count < 3`; Complete three times updating `count`; process completes.

**Acceptance Scenarios**:

1. **Given** a waiting-task activity with standard loop `testBefore=false` and condition that stays true for two iterations, **When** Complete runs twice with variables that satisfy then fail the condition, **Then** the activity ACTIVATED twice then the process continues.
2. **Given** `loopMaximum="2"` and a condition that would allow more, **When** Completes run, **Then** the activity executes at most twice then leaves.

---

### User Story 2 - Test-before skip (Priority: P1)

With `testBefore="true"`, if the condition is already false on arrival, the activity does not wait; the token continues downstream.

**Acceptance Scenarios**:

1. **Given** `testBefore=true` and a false condition on CreateInstance variables, **When** the instance is created, **Then** the activity does not remain waiting and the process reaches end (or the next element).

---

### User Story 3 - Recover mid-loop (Priority: P2)

A wait on a later loop iteration survives Recover; Complete continues the same loop rules.

**Acceptance Scenarios**:

1. **Given** a wait on the second iteration, **When** Recover then Complete until exit, **Then** outcome matches continuous execution.

---

### User Story 4 - Deploy mutual exclusion (Priority: P2)

Standard loop and multi-instance on the same activity remain mutually exclusive at Deploy.

**Acceptance Scenarios**:

1. **Given** both `standardLoopCharacteristics` and `multiInstanceLoopCharacteristics` on one activity, **When** Deploy runs, **Then** Deploy fails with `UNSUPPORTED_ELEMENT`.
2. **Given** only `standardLoopCharacteristics`, **When** Deploy runs, **Then** Deploy succeeds.

---

### Edge Cases

- Empty `loopCondition` and no `loopMaximum`: execute exactly once.
- Empty `loopCondition` with `loopMaximum=N`: execute N times.
- `loopCounter` is available to the condition (1-based for the current/next evaluation per BPMN).
- Job-backed tasks (Service / Business Rule / Script) re-wait as a job each iteration.
- SubProcess / Call Activity standard loop may be deferred if not wired in this increment; waiting task types are in scope.
- Complex MI behavior remains a separate Planned item.

## Requirements *(mandatory)*

- **FR-001**: Deploy MUST accept `standardLoopCharacteristics` on supported waiting activities (User / Manual / abstract Task / Receive / Service / Business Rule / Script).
- **FR-002**: Deploy MUST reject combining standard and multi-instance loop characteristics on the same element.
- **FR-003**: With `testBefore=false`, the engine MUST execute the activity then evaluate `loopCondition` (and max) to decide whether to repeat.
- **FR-004**: With `testBefore=true`, the engine MUST evaluate before each execution; a false condition on first arrival MUST skip the activity wait.
- **FR-005**: `loopMaximum` when set MUST cap the number of executions.
- **FR-006**: Recover mid-iteration MUST preserve loop progress and allow Complete to continue.
- **FR-007**: Existing multi-instance and non-loop suites MUST remain passing.

## Success Criteria *(mandatory)*

- **SC-001**: Do-while fixture repeats the expected number of waits then completes the process.
- **SC-002**: Test-before false-on-arrival fixture does not leave a wait on the looped activity.
- **SC-003**: Recover mid-loop then finish matches continuous.
- **SC-004**: Processing regression suite remains green.

## Assumptions

- Condition language is the existing `processing/expr` evaluator (same as conditional flows/catches).
- `loopCounter` injection follows BPMN (increment before each condition evaluation).
- SubProcess / Call Activity / Send Task standard loop can follow in a later slice if not included here.
