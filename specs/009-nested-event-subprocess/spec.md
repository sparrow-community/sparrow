# Feature Specification: Nested Event Sub-Process

**Feature Branch**: `009-nested-event-subprocess`

**Created**: 2026-09-08

**Status**: Draft

**Input**: User description: "引擎完整性优先；下一步 Remaining P3：Nested Event Sub-Process — 允许 Event Sub-Process 嵌套在另一个 Event Sub-Process 内：部署接受、外层触发后武装内层、内层可打断/处理外层作用域内事件，并覆盖 Recover。"

## User Scenarios & Testing *(mandatory)*

Event Sub-Processes already work at process level and nested inside **embedded** SubProcesses. Deploy today rejects an Event Sub-Process whose parent is another Event Sub-Process. This increment allows that nesting: when the outer Event Sub-Process is active, inner Event Sub-Processes arm for that scope; triggers and interrupting semantics follow existing Event Sub-Process rules with the outer ESP as parent scope.

### User Story 1 - Deploy and run ESP nested in ESP (Priority: P1)

An author places an Event Sub-Process inside another Event Sub-Process (e.g. outer message handler contains an inner message Event Sub-Process). Deploy succeeds. After the outer Event Sub-Process is triggered and running, the inner one is armed. Publishing the inner message interrupts (or handles) work inside the outer Event Sub-Process according to the inner start’s interrupting flag.

**Why this priority**: Namesake gap; unblocks the deploy reject.

**Independent Test**: Fixture with process-level outer ESP + inner ESP; trigger outer → waiting inside outer with inner armed → trigger inner → outer work terminated (if interrupting) and process completes via inner path.

**Acceptance Scenarios**:

1. **Given** an Event Sub-Process nested inside another Event Sub-Process, **When** Deploy runs, **Then** deploy succeeds.
2. **Given** the outer Event Sub-Process has been triggered and is waiting on an inner user task, **When** the instance is inspected, **Then** the nested Event Sub-Process is armed.
3. **Given** the nested Event Sub-Process is armed and interrupting, **When** its start message (or signal/timer) fires, **Then** active tokens in the outer Event Sub-Process are terminated and the nested Event Sub-Process runs to completion.

---

### User Story 2 - Disarm nested ESP when outer scope ends (Priority: P1)

When the outer Event Sub-Process completes without the inner firing, the nested Event Sub-Process is disarmed (same as ESP nested in embedded SubProcess).

**Why this priority**: Prevents stale arms after the outer handler finishes.

**Independent Test**: Trigger outer → complete outer path without inner event → nested arm gone; optional late inner message does not start the nested ESP.

**Acceptance Scenarios**:

1. **Given** outer Event Sub-Process completed successfully, **When** arms are inspected, **Then** the nested Event Sub-Process is not armed.
2. **Given** that state, **When** the nested start message is published, **Then** delivery count is 0 for that instance’s nested ESP (no new nested run).

---

### User Story 3 - Recover with nested ESP armed (Priority: P2)

Recover while waiting inside the outer Event Sub-Process with the nested Event Sub-Process armed; after Recover, the nested trigger still works.

**Why this priority**: Completeness Recover coverage.

**Independent Test**: Open/Recover mid-outer-wait; publish inner message → same outcome as continuous run.

**Acceptance Scenarios**:

1. **Given** outer ESP active and nested ESP armed, **When** Recover finishes, **Then** nested ESP remains armed.
2. **Given** Recover restored that state, **When** the nested start event fires, **Then** the process completes as in Story 1.

---

### Edge Cases

- ESP nested in **embedded** SubProcess remains unchanged (`m4_nested_event_subprocess`).
- Deeper than one ESP-in-ESP level: allowed if each nested ESP validates; arms when its parent ESP scope is entered.
- Non-interrupting nested ESP: SHOULD be supported if outer non-interrupting ESP already is (same runtime path); minimum coverage is interrupting nested.
- Call Activity / Transaction: out of scope.
- Multiple same-kind boundaries: separate Remaining item.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Deploy MUST accept an Event Sub-Process whose parent scope is another Event Sub-Process (remove the nested-in-ESP reject).
- **FR-002**: When an Event Sub-Process is entered (triggered), the engine MUST arm Event Sub-Processes whose parent scope is that Event Sub-Process id.
- **FR-003**: Nested Event Sub-Process trigger, interrupt, and complete semantics MUST match existing Event Sub-Process behavior with parent scope = outer Event Sub-Process.
- **FR-004**: Completing or terminating the outer Event Sub-Process scope MUST disarm nested Event Sub-Processes in that scope.
- **FR-005**: Event Sub-Process nested in embedded SubProcess and process-level Event Sub-Process behavior MUST remain unchanged.
- **FR-006**: Recover MUST restore nested Event Sub-Process arms while the outer Event Sub-Process is active so a subsequent trigger yields the same outcome.

### Key Entities

- **Outer Event Sub-Process**: triggeredByEvent SubProcess that contains another Event Sub-Process.
- **Nested Event Sub-Process**: Event Sub-Process with ParentScopeID equal to the outer Event Sub-Process id.
- **Arm**: Existing EventSubProcesses projection entry for a start catch.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Deploy of ESP-in-ESP fixture succeeds.
- **SC-002**: After outer trigger, nested ESP is armed; nested interrupting trigger terminates outer waiting work and completes via nested path.
- **SC-003**: Completing outer without nested fire disarms nested ESP.
- **SC-004**: Recover mid-outer with nested armed, then nested fire, matches continuous outcome.
- **SC-005**: Existing `m3_*` / `m4_nested_event_subprocess` / error ESP suites still pass.

## Assumptions

- Indexing already assigns ParentScopeID when walking nested FlowElements; the main blockers are deploy validation and arming on ESP enter.
- Supported start kinds for nested ESP are the same as today (message, signal, timer, error).
- Multiple same-kind boundaries and new element types stay out of scope.
