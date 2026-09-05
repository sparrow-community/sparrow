# Feature Specification: Multi-Instance Call Activity

**Feature Branch**: `006-multi-instance-callactivity`

**Created**: 2026-09-05

**Status**: Draft

**Input**: User description: "引擎完整性优先；下一步按 Remaining P2：在已有 Multi-Instance（User/Service/SubProcess）与 Call Activity（含跨部署、boundary/compensation）之上，补齐 Multi-Instance Call Activity，使一次调用可为每次循环启动独立子流程实例，并覆盖并行/顺序、集合输入、完成条件、边界打断与 Recover。"

## User Scenarios & Testing *(mandatory)*

Authors already deploy Call Activities (same-file or cross-deployment), map IO, attach boundaries/compensation, and run multi-instance loops on User Task, Service Task, and embedded Sub-Process. This increment allows **loop characteristics on Call Activity** so each inner iteration starts its own called process instance; the parent continues only when the loop join rules say so.

### User Story 1 - Parallel multi-instance Call Activity (Priority: P1)

An author marks a Call Activity with parallel multi-instance and a fixed cardinality (for example `3`). When the parent reaches that Call Activity, three child process instances start (one per inner iteration). Completing all three children completes the Call Activity once and the parent continues on the single outgoing path.

**Why this priority**: Proves the core gap — multi-instance host plus one child instance per inner — without collection or early completion complexity.

**Independent Test**: Deploy caller + called processes; Call Activity with parallel MI cardinality 3. Start parent; three distinct child instances wait on the called user task. Complete all three children; parent ends once.

**Acceptance Scenarios**:

1. **Given** a parallel multi-instance Call Activity with cardinality 3, **When** the parent reaches it, **Then** three child process instances exist and the parent does not continue past the Call Activity until all three children complete.
2. **Given** two of three children have completed, **When** the third completes, **Then** the Call Activity completes exactly once and the outgoing sequence flow is taken once.
3. **Given** three children are active, **When** an operator inspects them, **Then** each has a distinct process instance id and each links back to the same parent Call Activity host.

---

### User Story 2 - Sequential multi-instance Call Activity (Priority: P1)

An author marks a Call Activity with sequential multi-instance and fixed cardinality. Only one child process instance is active at a time; completing it starts the next until all iterations finish.

**Why this priority**: Same join machinery as parallel with a concurrency constraint; common for ordered fan-out to callees.

**Independent Test**: Sequential MI Call Activity cardinality 3. At most one child waits at a time. After three child completions in order, the parent ends.

**Acceptance Scenarios**:

1. **Given** sequential multi-instance Call Activity with cardinality 3, **When** the parent arrives, **Then** exactly one child process instance is active.
2. **Given** the first child has completed, **When** the engine advances, **Then** a second child starts (the first does not reopen).
3. **Given** all three children complete in order, **When** the last finishes, **Then** the Call Activity completes once.

---

### User Story 3 - Collection input and IO mapping per child (Priority: P2)

An author binds multi-instance input to a collection and optional Call Activity IO mappings. The number of children equals the collection size. Each child receives the current element (and mapped inputs); optional output collection assembles mapped outputs from completed children.

**Why this priority**: Real “for each item, call this process” patterns need data-driven fan-out, not only fixed cardinality.

**Independent Test**: Parent variable `items` = `["a","b"]`. Parallel MI Call Activity over that collection with input mapping into the child. Two children start with the corresponding element; complete both; optional output collection is present on the parent if configured.

**Acceptance Scenarios**:

1. **Given** a collection of N elements, **When** the MI Call Activity activates, **Then** N child instances start.
2. **Given** an inner Call Activity iteration is active, **When** the child variables are inspected, **Then** mapped inputs / element variable reflect that iteration’s collection entry.
3. **Given** an empty collection, **When** the flow arrives, **Then** zero children start and the Call Activity completes immediately.

---

### User Story 4 - Completion condition and interrupting boundary (Priority: P2)

An author sets a completion condition so the Call Activity may finish before all children complete; remaining children are terminated. Separately, an interrupting boundary on the multi-instance Call Activity cancels the loop host and all active children, then takes the boundary path once.

**Why this priority**: Quorum / first-success and timeout/cancel on fan-out are required for parity with other MI hosts and existing Call Activity boundary behavior.

**Independent Test**: Parallel MI Call Activity cardinality 5 with completion condition “at least 2 completed”; after two children complete, remaining children are terminated and the parent continues. Separate fixture: interrupting timer on MI Call Activity; after fire, all children terminated and timeout path taken.

**Acceptance Scenarios**:

1. **Given** a completion condition requiring two of five children, **When** any two complete, **Then** the Call Activity completes and still-active children are terminated.
2. **Given** an interrupting boundary on a multi-instance Call Activity with active children, **When** the boundary fires, **Then** all active children and the loop host are terminated and the boundary path is taken once.
3. **Given** no completion condition, **When** the Call Activity runs, **Then** default `All` behavior applies (all started children must complete).

---

### User Story 5 - Recover mid multi-instance call (Priority: P2)

An operator recovers the engine while a multi-instance Call Activity has some children completed and others still waiting. After recover, remaining children can be completed and the parent reaches the same outcome as without a crash.

**Why this priority**: Completeness requires Recover coverage for unfinished COMMAND chains and rebuilt loop/child linkage.

**Independent Test**: Parallel MI Call Activity cardinality 3; complete one child; recover from the event log; complete the other two; parent completes. Child deployment ids and parent/child links remain correct (including cross-deployment callees).

**Acceptance Scenarios**:

1. **Given** one of three children completed before crash, **When** Recover finishes, **Then** two waiting children (or equivalent waiting work) remain and loop counters match.
2. **Given** Recover restored mid-loop state, **When** the remaining children complete, **Then** the Call Activity and parent complete once.
3. **Given** a cross-deployment callee, **When** Recover runs mid-loop, **Then** each child still binds to the callee deployment used at start (not a newer revision).

---

### Edge Cases

- What happens when `loopCardinality` is zero or negative? Zero children; Call Activity completes immediately.
- What happens when the collection variable is missing or not an array? Treat as empty (immediate completion), same as other multi-instance hosts.
- What happens to compensation on the Call Activity host? Subscription occurs when the multi-instance Call Activity host completes successfully (not per unfinished child); interrupted / early-cancelled loops do not subscribe for unfinished work, consistent with other activities.
- How does deploy treat multi-instance on Call Activity that was previously unsupported? Deploy MUST accept valid `multiInstanceLoopCharacteristics` on Call Activity once this increment ships; invalid cardinality/collection config still rejects.
- What happens if one child ends with an uncaught error? Existing error propagation for Call Activity / child termination applies to that iteration; join and completion-condition rules then apply as for other MI hosts.
- Non-multi-instance Call Activity (same-file and cross-deploy) and multi-instance on User/Service/SubProcess MUST remain unchanged.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Deploy MUST accept `multiInstanceLoopCharacteristics` on Call Activity when the element is otherwise a valid Call Activity (bundled or external callee).
- **FR-002**: Parallel multi-instance Call Activity MUST allow multiple child process instances active concurrently under one parent loop host.
- **FR-003**: Sequential multi-instance Call Activity MUST allow at most one child process instance active at a time.
- **FR-004**: Each inner iteration MUST start exactly one called process instance (same-file or cross-deployment resolution rules unchanged per child).
- **FR-005**: Loop cardinality or collection input MUST determine the number of children using the same rules as multi-instance User/Service/SubProcess.
- **FR-006**: Default completion MUST require all started children to complete (`All`) when no completion condition is declared.
- **FR-007**: A declared completion condition MUST be evaluated after each child completes or is cancelled and MAY complete the Call Activity early when true.
- **FR-008**: Early completion or interrupting boundary on the multi-instance Call Activity MUST terminate still-active children consistently with existing terminate-child behavior.
- **FR-009**: After the multi-instance Call Activity completes, the outgoing sequence flow MUST be taken exactly once.
- **FR-010**: Call Activity IO mapping MUST apply per child (inputs at child start; outputs contributed per completed iteration; optional output collection when configured).
- **FR-011**: Recover MUST restore loop counters, host/inner linkage, and parent/child instance relationships so remaining work yields the same outcome as before the crash.
- **FR-012**: Existing non-multi-instance Call Activity and multi-instance behavior on other element types MUST remain unchanged.
- **FR-013**: Boundaries and compensation on multi-instance Call Activity MUST follow the same host-level rules as other multi-instance activities (cancel all inners on interrupt; compensate subscribe on successful host complete).

### Key Entities

- **Multi-instance Call Activity host**: Parent-side loop host for the Call Activity element.
- **Inner call iteration**: One loop index that owns one called process instance and optional element input/output.
- **Called process instance**: Child instance started for an inner iteration; linked to parent instance and host token/iteration.
- **Loop counters**: Total, active, and completed iteration counts for join and completion conditions.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Parallel MI Call Activity with cardinality 3 starts three distinct child instances before the parent can continue.
- **SC-002**: Sequential MI Call Activity never has more than one active child at a time.
- **SC-003**: Collection size N yields N children; empty collection completes with zero children.
- **SC-004**: Completion condition “2 of 5” continues the parent after two child completions with no stray active children.
- **SC-005**: Interrupting boundary on MI Call Activity terminates all active children and takes the boundary path once.
- **SC-006**: Recover mid-loop (1 of 3 children done) then finishing the rest completes the parent identically to a no-crash run.
- **SC-007**: Regression: existing Call Activity and multi-instance suites continue to pass.

## Assumptions

- Baseline includes multi-instance on User Task, Service Task, SubProcess (`002`) and Call Activity with cross-deployment, boundary, and compensation (`004`, `005`).
- Loop cardinality, collection input, completion condition, and loop metadata semantics reuse the existing multi-instance rules; this increment extends the host type set to Call Activity.
- Cross-deployment resolution remains latest revision at each child start unless a prior increment already pinned per publication; Recover must not rebind a child to a newer callee revision mid-redrive.
- Complex behavior definitions / none-one behavior event refs stay out of scope (same as `002`).
- Instantiate event-based gateway, nested Event Sub-Process, live instance migration, and other Remaining/Future items stay out of scope.
- Compensation into unfinished children or unfinished SubProcess remains a separate Remaining item.
