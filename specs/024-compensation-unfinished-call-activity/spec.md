# Feature Specification: Compensation into Unfinished Call Activity Child

**Feature Branch**: `024-compensation-unfinished-call-activity`

**Created**: 2026-09-15

**Status**: Draft

**Input**: User description: "在本地继续推进 Planned：Compensation into an unfinished Call Activity child instance。说明与文档中不要使用 MI、EBG、ESB 等简写。"

## User Scenarios & Testing *(mandatory)*

Spec 008 taught parent compensate throws to enter unfinished embedded SubProcesses. Call Activity was left out: compensation stopped at the Call Activity host, and only a Call Activity compensation boundary after the child completed could undo the call. This increment closes that gap: when a parent-scope compensate throw fires while a Call Activity host is still waiting on an active child instance, the engine cancels unfinished work in the child, runs compensation handlers for activities that already completed inside the child (on the child instance), terminates the child and the Call Activity host, then continues the parent throw.

### User Story 1 - Broadcast compensate enters unfinished Call Activity child (Priority: P1)

An author models a parallel split: one branch starts a Call Activity whose child completes Task_A (with a compensation handler) then waits on Task_B; the sibling branch reaches a compensate throw. When the throw fires, Task_B is terminated, Undo_A runs on the child, the child and Call Activity host end, and the parent throw path can finish.

**Independent Test**: Parallel fixture → complete Task_A in child → leave Task_B waiting → complete sibling → Undo_A waiting on child → Complete Undo_A → process completed.

**Acceptance Scenarios**:

1. **Given** an active Call Activity with completed compensatable Task_A and waiting Task_B in the child, **When** a parent-scope broadcast compensate throw fires, **Then** Task_B is terminated and Undo_A becomes waiting on the child instance.
2. **Given** that situation, **When** Undo_A is Completed, **Then** the child instance terminates, the Call Activity host is terminated, and the parent compensate throw completes so the caller process can complete.

---

### User Story 2 - Targeted activityRef to unfinished Call Activity (Priority: P1)

A compensate throw with `activityRef` equal to the unfinished Call Activity id enters only that child. Unrelated parent-scope compensation subscriptions are not consumed by that throw.

**Acceptance Scenarios**:

1. **Given** unfinished CallActivity_1 and a completed sibling Task_X both with handlers, **When** compensate throw has `activityRef=CallActivity_1`, **Then** the child is cancelled and inner Undo runs; Task_X’s subscription remains.
2. **Given** `activityRef` names a completed activity (not an active Call Activity), **When** compensate throws, **Then** existing targeted behavior is unchanged.

---

### User Story 3 - Recover mid child compensation (Priority: P2)

An operator recovers while Undo_A is waiting on the child after the parent throw entered the unfinished Call Activity. After Recover, completing Undo_A finishes the chain with the same outcome as a continuous run.

**Acceptance Scenarios**:

1. **Given** parent throw waiting and child Undo_A waiting, **When** Recover finishes, **Then** both waits are restored.
2. **Given** Recover restored that state, **When** Undo_A is Completed, **Then** caller and child reach the same terminal outcome as without a crash.

---

### Edge Cases

- Unfinished Call Activity whose child has **no** completed compensatable activities: terminate child work and Call Activity host; parent throw continues with an empty contribution from that call.
- Broadcast MUST also keep unfinished embedded SubProcess entry (spec 008) and same-scope completed subscriptions.
- Completed Call Activity compensation via boundary + association remains unchanged (`m9_call_compensate`).
- Nested Call Activity inside the child: terminating the child instance terminates nested children as today; recursive “unfinished into unfinished Call Activity” as a separate parent throw is not required beyond terminateInstance’s existing nested terminate.
- Cross-deployment called process: same behavior when callee deployment differs.
- Documentation spells out “Call Activity” and “event sub-process” without abbreviation.

## Requirements *(mandatory)*

- **FR-001**: When a compensate throw runs in scope S and a Call Activity host in S is still waiting with an active child instance, the engine MUST enter that child for unfinished-call compensation (cancel active child work, run child CompensationSubs in reverse order on the child instance).
- **FR-002**: After child compensation handlers finish (or the child had none), the engine MUST terminate the child instance and the Call Activity host (without completing the Call Activity successfully).
- **FR-003**: Broadcast compensate MUST apply FR-001/FR-002 for each unfinished direct-child Call Activity under the throw scope, in addition to unfinished SubProcess entry and same-scope subscriptions from spec 008.
- **FR-004**: Compensate with `activityRef` equal to an unfinished Call Activity id MUST apply FR-001/FR-002 for that Call Activity only and MUST NOT consume unrelated same-scope subscriptions.
- **FR-005**: Child compensation handlers MUST execute on the child instance with the callee deployment (not on the caller).
- **FR-006**: The parent PendingCompensation MUST wait until unfinished-call child compensation finishes before advancing remaining parent handler queue / completing the throw.
- **FR-007**: Recover MUST restore parent throw wait, Call Activity host link, and child handler wait so Completing the handler yields the same outcome.
- **FR-008**: Existing suites (`m3_compensation`, `m4_compensate_*`, `m9_call_compensate`, unfinished SubProcess, compensation event sub-process) MUST remain green.

## Success Criteria *(mandatory)*

- **SC-001**: Parallel unfinished Call Activity fixture: parent compensate terminates child Task_B and runs Undo_A on the child; caller completes after Undo_A.
- **SC-002**: Targeted `activityRef` to unfinished Call Activity does not consume an unrelated parent-scope subscription.
- **SC-003**: Recover mid Undo_A then Complete matches continuous-run finals.
- **SC-004**: `m9_call_compensate` and unfinished SubProcess suites still pass.
- **SC-005**: Processing regression suite remains green.

## Assumptions

- Reuse `PublicationResumeParent` with `Completed=false` to terminate the Call Activity host after child compensation.
- Parent `PendingCompensation` gains a waiting-children counter (or equivalent) so the throw does not complete early.
- Condition and message languages unchanged.
