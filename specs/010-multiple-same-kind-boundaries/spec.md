# Feature Specification: Multiple Same-Kind Boundaries

**Feature Branch**: `010-multiple-same-kind-boundaries`

**Created**: 2026-09-08

**Status**: Draft

**Input**: User description: "引擎完整性优先；下一步 Remaining P3：Multiple same-kind boundaries on one activity — 允许同一活动挂多个同类型边界（如多个不同 message/timer/signal），部署与运行时均武装并正确触发/取消；Recover 覆盖。"

## User Scenarios & Testing *(mandatory)*

Today an activity may have at most one timer, one message, and one signal boundary (different kinds together are OK). Deploy rejects a second boundary of the same kind. Authors need multiple message (or timer/signal) boundaries on one activity—e.g. cancel vs escalate messages, or short vs long timers.

### User Story 1 - Multiple message boundaries (Priority: P1)

An author attaches two or more message boundaries with **different** message names to one User Task. Deploy succeeds. While the task waits, all message subscriptions are armed. Publishing one interrupting message completes that boundary path, terminates the host and sibling boundaries, and follows that boundary’s outgoing flow.

**Why this priority**: Most common same-kind case (Camunda: unique message names per activity).

**Independent Test**: Fixture with three message boundaries on one task; publish one → that path completes; siblings terminated.

**Acceptance Scenarios**:

1. **Given** an activity with two+ message boundaries with distinct names, **When** Deploy runs, **Then** deploy succeeds.
2. **Given** the activity is waiting, **When** one of those messages is published (interrupting), **Then** that boundary COMPLETED, host TERMINATED, other message boundaries TERMINATED, process continues on that boundary’s path.
3. **Given** duplicate message names on two boundaries of the same activity, **When** Deploy runs, **Then** deploy fails with UNSUPPORTED_ELEMENT.

---

### User Story 2 - Multiple timer boundaries (Priority: P1)

An author attaches two timer boundaries with different durations to one activity. Both are armed. The first due interrupting timer wins; the later timer is terminated without taking its path.

**Why this priority**: Same-kind timer races are a standard pattern.

**Independent Test**: PT0S and PT1H timers; FireDue → short path wins; long timer TERMINATED.

**Acceptance Scenarios**:

1. **Given** two timer boundaries on one waiting activity, **When** the earlier timer is due and FireDue runs, **Then** that boundary path is taken and the other timer is terminated.

---

### User Story 3 - Recover with multiple same-kind waits (Priority: P2)

Recover while an activity has multiple message (or timer) boundaries armed; after Recover, publishing/firing still selects the correct boundary.

**Acceptance Scenarios**:

1. **Given** two message boundaries armed, **When** Recover runs, **Then** both remain armed.
2. **Given** Recover restored arms, **When** one message is published, **Then** outcome matches the continuous run.

---

### Edge Cases

- Multiple signal boundaries with distinct signal names: same rules as messages (in scope if time permits; minimum is message + timer).
- Compensation: still at most one compensation boundary per activity.
- Error: multiple error boundaries with different codes remain allowed (existing); same code still rejected.
- Mixing kinds (timer+message+signal) unchanged.
- Non-interrupting multiples: first increment MAY focus interrupting; non-interrupting multi-message SHOULD work if collectors and cancel-one paths allow.

## Requirements *(mandatory)*

- **FR-001**: Deploy MUST allow multiple timer boundaries on one activity.
- **FR-002**: Deploy MUST allow multiple message boundaries on one activity when message names differ; MUST reject duplicate message names on the same activity.
- **FR-003**: Deploy MUST allow multiple signal boundaries when signal names differ; MUST reject duplicates; MUST still reject a second compensation boundary on the same activity.
- **FR-004**: On activity enter, the engine MUST arm all attached timer/message/signal boundaries for that activity.
- **FR-005**: Triggering one interrupting boundary MUST terminate the host and all sibling waiting boundaries on that activity.
- **FR-006**: Existing single-boundary and mixed-kind behavior MUST remain unchanged.
- **FR-007**: Recover MUST restore all armed same-kind boundaries so a subsequent trigger matches a continuous run.

## Success Criteria *(mandatory)*

- **SC-001**: Deploy accepts ≥2 distinct message boundaries and ≥2 timer boundaries on one activity.
- **SC-002**: Publishing one of three interrupting messages completes only that path; siblings terminated.
- **SC-003**: Earlier of two timers wins via FireDue.
- **SC-004**: Recover mid-wait with two message arms then publish matches continuous outcome.
- **SC-005**: Existing boundary suites still pass.

## Assumptions

- Unique message/signal names per activity match Camunda guidance.
- Proto may gain a repeated waiting-boundary payload; regenerate via `build.sh`.
- Escalation and other Remaining element types stay out of scope.
