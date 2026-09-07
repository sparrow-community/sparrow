# Feature Specification: Compensation into Unfinished SubProcess

**Feature Branch**: `008-compensation-unfinished-subprocess`

**Created**: 2026-09-07

**Status**: Draft

**Input**: User description: "引擎完整性优先；下一步 Remaining P3：Compensation into unfinished SubProcess — 当父作用域抛出补偿时，若嵌入式 SubProcess 仍在执行，取消该子流程中未完成工作，并对已在该子流程内成功完成且已订阅补偿的活动按逆序执行补偿处理器；覆盖 Recover。"

## User Scenarios & Testing *(mandatory)*

Today compensation runs only for subscriptions whose activity lives in the **same scope** as the compensate throw, and only after those activities have completed. If an embedded SubProcess is still active when the parent throws compensation, completed work **inside** that SubProcess keeps its subscriptions but those handlers are never invoked (matching common Camunda “do not propagate into unfinished subprocess” default). This increment adds the missing completeness behavior: parent-scope compensation **enters** unfinished embedded SubProcesses — terminate active tokens in that SubProcess, then run compensation handlers for activities that already completed inside it.

### User Story 1 - Parent compensate cancels unfinished SubProcess and undoes completed inners (Priority: P1)

An author models a parallel split: one branch runs an embedded SubProcess (inner activity A completes and has a compensation handler; then inner activity B waits); the other branch completes some work and reaches a compensate throw (broadcast or targeting the SubProcess). When the throw fires while the SubProcess is still waiting on B, the engine terminates unfinished work inside the SubProcess and runs A’s compensation handler before the throw path continues.

**Why this priority**: Namesake gap; proves terminate + nested-scope compensation queue.

**Independent Test**: Parallel fixture → complete A inside SP → leave B waiting → complete sibling path to compensate throw → Undo_A runs; B terminated; process can complete.

**Acceptance Scenarios**:

1. **Given** an active SubProcess with completed compensatable A and waiting B, **When** a parent-scope compensate throw fires, **Then** tokens inside the SubProcess (including B and the SubProcess host as needed) are terminated.
2. **Given** that situation, **When** compensation runs, **Then** A’s compensation handler becomes waiting (or runs to completion if auto), and completing it allows the compensate throw path to finish.
3. **Given** the same model but the SubProcess already fully completed before the throw, **When** compensate fires, **Then** existing completed-SubProcess compensation behavior remains (SubProcess-level handler and/or already-subscribed same-scope rules — no regression on `m4_compensate_subprocess`).

---

### User Story 2 - Targeted compensate by activityRef to unfinished SubProcess (Priority: P1)

An author sets `activityRef` on the compensate throw to the unfinished SubProcess id. Only that SubProcess is entered for cancel + inner compensation; other parent-scope subscriptions are not consumed by that targeted throw (same targeting rules as today for completed activities, extended to unfinished SubProcess entry).

**Why this priority**: Matches BPMN activityRef targeting; safer than always broadcasting into every active SP.

**Independent Test**: Fixture with activityRef=SubProcess_1 while SP active → inner undo runs; unrelated sibling completed activity’s handler is not run by this throw.

**Acceptance Scenarios**:

1. **Given** unfinished SubProcess_1 and a completed sibling Task_X both with handlers, **When** compensate throw has `activityRef=SubProcess_1`, **Then** SubProcess_1 is cancelled and inner completed handlers run; Task_X’s subscription remains.
2. **Given** `activityRef` names a completed activity (not an active SubProcess), **When** compensate throws, **Then** existing targeted behavior is unchanged.

---

### User Story 3 - Recover mid unfinished-SubProcess compensation (Priority: P2)

An operator recovers while a parent compensate throw is waiting on an inner compensation handler after cancelling an unfinished SubProcess. After Recover, completing the handler finishes the throw path with the same outcome.

**Why this priority**: Completeness requires Recover coverage.

**Independent Test**: Drive to Undo waiting after SP cancel; Open/Recover; Complete Undo → process completes like continuous run.

**Acceptance Scenarios**:

1. **Given** PendingCompensation after unfinished-SP entry with a waiting handler, **When** Recover finishes, **Then** the handler is still waiting and throw is still pending.
2. **Given** Recover restored that state, **When** the handler is completed, **Then** the process reaches the same completed outcome as without a crash.

---

### Edge Cases

- What if the unfinished SubProcess has **no** completed compensatable activities inside? Terminate active tokens; compensate throw continues with an empty inner queue (no handler).
- What if nested SubProcess (SP2 inside SP1) is unfinished when parent throws? This increment MUST handle **one** level of unfinished embedded SubProcess under the throw scope; deeper unfinished nesting MAY be terminated with the outer unfinished SP and only run subscriptions whose activity scope is that unfinished SP (not recursive unfinished-into-unfinished beyond one documented level — prefer: terminate all descendant tokens under the unfinished SP and run all CompensationSubs whose activity ScopeOf is that SP id or a descendant scope under it).
- Broadcast (empty activityRef): MUST enter **each** unfinished embedded SubProcess that is a direct child scope of the throw scope (cancel + inner handlers), in addition to running same-scope completed subscriptions as today.
- Call Activity unfinished child: out of scope (compensation still stops at Call Activity; use Call Activity compensation boundary after complete).
- Transaction SubProcess / cancel events: out of scope.
- Compensation Event Sub-Process as handler: out of scope.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: When a compensate throw runs in scope S and an embedded SubProcess P is still active with `ScopeOf(P.host)=S` (or P is a direct child scope of S), the engine MUST terminate active tokens belonging to P (and descendants) before or as part of starting compensation for that throw’s unfinished-SP entry.
- **FR-002**: After terminating unfinished P, the engine MUST queue compensation handlers for subscriptions whose compensated activity’s scope is P (completed activities inside P), in reverse completion order, as part of that compensate throw’s handler queue.
- **FR-003**: Broadcast compensate (no activityRef) MUST still run same-scope completed subscriptions as today, and MUST also apply FR-001/FR-002 for each unfinished direct-child embedded SubProcess under the throw scope.
- **FR-004**: Compensate with `activityRef` equal to an unfinished SubProcess id MUST apply FR-001/FR-002 for that SubProcess only and MUST NOT consume unrelated same-scope subscriptions.
- **FR-005**: Existing compensation for completed activities / completed SubProcess boundary handlers MUST remain unchanged when no unfinished SubProcess entry applies.
- **FR-006**: Recover MUST restore PendingCompensation and waiting handlers after unfinished-SubProcess compensation has started so completing remaining handlers yields the same outcome.
- **FR-007**: Call Activity, Transaction/cancel, and compensation Event Sub-Process handlers stay out of scope for this increment.

### Key Entities

- **Unfinished SubProcess**: Embedded SubProcess host still active (waiting or running tokens inside) when parent compensate throw starts.
- **Inner compensation subscription**: CompensationSub for an activity whose ScopeOf is the SubProcess id.
- **PendingCompensation**: Existing throw wait state; queue may include handlers discovered via unfinished-SP entry.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Parallel unfinished-SP fixture: parent compensate terminates waiting inner work and runs the completed inner’s undo handler.
- **SC-002**: Targeted `activityRef` to unfinished SubProcess does not consume an unrelated parent-scope subscription.
- **SC-003**: `m3_compensation` / `m4_compensate_subprocess` / Call Activity compensate suites still pass.
- **SC-004**: Recover after SP cancel + undo waiting, then complete undo → same final status/intents as continuous run.
- **SC-005**: Unfinished SP with zero inner subscriptions: SP work terminated; throw completes without requiring a handler.

## Assumptions

- Embedded SubProcess termination helpers already exist for error / Event Sub-Process paths and can be reused conceptually.
- Compensation handlers remain userTask/serviceTask `isForCompensation` as today.
- Parallel gateway (or equivalent concurrent paths) is an acceptable way to keep a SubProcess unfinished while another path throws compensate.
- Deeper than one unfinished nesting is handled by terminating the whole unfinished subtree and selecting subscriptions by scope under that SubProcess; recursive “unfinished into unfinished” as separate throws is not required.
- Nested Event Sub-Process and multiple same-kind boundaries remain separate Remaining items.
