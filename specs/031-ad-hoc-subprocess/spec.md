# Feature Specification: Ad-Hoc SubProcess

**Feature Branch**: `main` (developed on main)

**Created**: 2026-09-16

**Status**: Draft

**Input**: User description: "继续推进 Planned：Ad-Hoc SubProcess。"

## User Scenarios & Testing *(mandatory)*

An Ad-Hoc SubProcess holds inner activities with no sequence flow between them. Enabled inner activities are performed by their normal Complete / Job path in any order; the scope finishes when `completionCondition` becomes true (or when no enabled inner activity remains). `ordering` selects whether all inner activities are enabled at once (`Parallel`, default) or one at a time in document order (`Sequential`). `cancelRemainingInstances` decides whether a true completion condition cancels inner activities still running.

### User Story 1 - Parallel ordering with completion condition (Priority: P1)

An Ad-Hoc SubProcess with three inner User Tasks enables all three. Completing two of them makes `${done >= 2}` true, the third is cancelled, and the scope continues on its outgoing flow.

### User Story 2 - Sequential ordering (Priority: P1)

The same scope with `ordering="Sequential"` enables only the first inner activity. Each completion enables the next inner activity in document order until the completion condition is true.

### User Story 3 - Keep remaining instances (Priority: P2)

With `cancelRemainingInstances="false"`, a true completion condition does not cancel inner activities still running; the scope completes after the last enabled inner activity finishes.

### User Story 4 - Exhaustion and Recover (Priority: P2)

If every enabled inner activity finishes while the completion condition is still false, the scope completes instead of deadlocking. Recover in the middle of an Ad-Hoc SubProcess restores the parked host token and every enabled inner activity, and the scope still completes correctly afterwards.

## Requirements *(mandatory)*

- **FR-001**: Deploy MUST accept `adHocSubProcess` as `TYPE_AD_HOC_SUB_PROCESS`, indexing its inner activities in the scope of the Ad-Hoc SubProcess.
- **FR-002**: Deploy MUST require a non-empty `completionCondition` and MUST reject `ordering` values other than `Parallel` / `Sequential` with `UNSUPPORTED_ELEMENT`.
- **FR-003**: Deploy MUST reject an Ad-Hoc SubProcess whose body is not a flat set of supported waiting / job activities (no sequence flow, no events, no gateways, no nested scope, no boundary event, no multi-instance or standard loop inside), and MUST reject multi-instance / standard loop / `triggeredByEvent` on the Ad-Hoc SubProcess itself.
- **FR-004**: Entering the Ad-Hoc SubProcess MUST park a host token on the scope and enable inner activities: all of them for `Parallel` (default), the first one in document order for `Sequential`.
- **FR-005**: Completing an enabled inner activity MUST re-evaluate `completionCondition` against instance variables; inner activities have no outgoing sequence flow and MUST NOT be treated as a flow continuation.
- **FR-006**: When `completionCondition` is true, remaining enabled inner activities MUST be terminated when `cancelRemainingInstances` is absent or true; when it is false the scope MUST wait for them and complete after the last one finishes.
- **FR-007**: When `completionCondition` is false, `Sequential` MUST enable the next inner activity in document order and `Parallel` MUST keep waiting for the remaining ones; when no enabled inner activity remains the scope MUST complete.
- **FR-008**: Scope completion MUST emit the Ad-Hoc SubProcess COMPLETING / COMPLETED pair on the host token and then leave via its outgoing flow(s), with compensation subscription and attached-boundary cancellation handled like other activities.
- **FR-009**: Recover MUST rebuild the parked host token and all enabled inner activity tokens from the ledger.

## Success Criteria *(mandatory)*

- **SC-001**: Parallel ordering enables all inner activities; a true completion condition cancels the rest and the scope continues.
- **SC-002**: Sequential ordering enables exactly one inner activity at a time in document order.
- **SC-003**: `cancelRemainingInstances="false"` completes only after the last enabled inner activity finishes.
- **SC-004**: Exhaustion with a false condition completes the scope; Recover mid-scope preserves host and inner tokens.
- **SC-005**: Deploy rejects missing `completionCondition`, bad `ordering`, and non-flat Ad-Hoc bodies with `UNSUPPORTED_ELEMENT`.
