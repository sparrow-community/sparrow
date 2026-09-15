# Feature Specification: Abstract Task

**Feature Branch**: `019-abstract-task`

**Created**: 2026-09-15

**Status**: Draft

**Input**: User description: "继续 Planned P1：Abstract bpmn:task as Manual-equivalent wait/complete。Deploy 接受抽象 task；运行时等待与 Complete 与 Manual Task 相同；事件类型为 TYPE_TASK；支持边界与多实例。"

## User Scenarios & Testing *(mandatory)*

Authors often place an untyped `task` in executable diagrams. Deploy currently rejects it. This increment accepts abstract `task` with the same wait → Complete semantics as Manual Task, while recording `TYPE_TASK` on the ledger.

### User Story 1 - Abstract task waits and completes (Priority: P1)

An author models start → abstract `task` → end. CreateInstance leaves a wait on that task. Complete advances the token and the process finishes. Events show `TYPE_TASK` (not Manual).

**Why this priority**: Closes the last basic task-type gap called out in Planned P1.

**Independent Test**: Deploy fixture with `<task>`; CreateInstance → Complete → process completed with TASK events.

**Acceptance Scenarios**:

1. **Given** an executable process with an abstract `task`, **When** Deploy runs, **Then** Deploy succeeds.
2. **Given** that deployment, **When** CreateInstance runs, **Then** the instance waits on the abstract task.
3. **Given** that wait, **When** Complete runs for that element/token, **Then** the process completes and events include `TYPE_TASK` ACTIVATED and COMPLETED.

---

### User Story 2 - Recover mid-wait (Priority: P2)

A wait on abstract task survives Recover and Complete finishes the same way as continuous execution.

**Acceptance Scenarios**:

1. **Given** a waiting abstract task, **When** Recover then Complete, **Then** the process completes like the continuous path.

---

### User Story 3 - Boundaries and multi-instance (Priority: P2)

Abstract task may host interrupting boundaries and multi-instance like Manual Task.

**Acceptance Scenarios**:

1. **Given** abstract task with a message boundary, **When** PublishMessage fires while waiting, **Then** interrupting boundary path works.
2. **Given** multiInstanceLoopCharacteristics on abstract task, **When** a collection is provided, **Then** loop host/inner waits follow existing Manual MI rules.

---

### Edge Cases

- Abstract `task` is not a Job; Activate does not claim it.
- `manualTask` remains a distinct type (`TYPE_MANUAL_TASK`); no aliasing of types.
- Global tasks / non-process `task` constructs remain out of scope.
- Collaboration and choreography remain Excluded.

## Requirements *(mandatory)*

- **FR-001**: Deploy MUST accept abstract `bpmn:task` and index `TYPE_TASK`.
- **FR-002**: Abstract task MUST wait for Complete with the same token semantics as Manual Task.
- **FR-003**: Ledger Element intents for abstract task MUST use `TYPE_TASK`.
- **FR-004**: Complete / Recover MUST match Manual Task wait semantics.
- **FR-005**: Boundaries and multi-instance MUST be supported like Manual Task.
- **FR-006**: Existing suites MUST remain passing (the prior “abstract task rejected” expectation MUST be updated).

## Success Criteria *(mandatory)*

- **SC-001**: Abstract-task fixture Deploy → CreateInstance → Complete → process completed with TASK events.
- **SC-002**: Recover mid-wait then Complete matches continuous.
- **SC-003**: Processing regression suite remains green.

## Assumptions

- Semantics match Manual Task (external human/agent Complete); not Service/Script job-backed.
- Proto already defines `TYPE_TASK`.
- BPMN parse already maps `<task>` to `FlowElements.Tasks`.
