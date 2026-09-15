# Feature Specification: Activity Parallel Fan-out

**Feature Branch**: `main` (developed on main)

**Created**: 2026-09-15

**Status**: Draft

**Input**: User description: "继续下一步 Planned：Parallel fan-out from an activity with multiple unconditional outgoings。保持在主分支。"

## User Scenarios & Testing *(mandatory)*

When an activity (or start event) has multiple outgoing sequence flows with no `conditionExpression` and no `default`, BPMN takes all of them (implicit parallel split). Spec 013 kept take-first for that case; this increment replaces take-first with parallel fan-out. Conditional / default selection from 013 stays exclusive.

### User Story 1 - Activity fans out on Complete (Priority: P1)

A User Task has two unconditional outgoings to two User Tasks that join at a parallel gateway. Completing the first task creates tokens on both branches; completing both joiners completes the process.

### User Story 2 - Start Event fans out (Priority: P1)

A none Start Event with two unconditional outgoings forks both branches on CreateInstance (same join pattern).

### User Story 3 - Conditional activity paths unchanged (Priority: P1)

Existing conditional + default activity fixtures (`m17_conditional_activity`) remain exclusive (one path).

### User Story 4 - Recover (Priority: P2)

Recover while waiting on the fan-out User Task; Complete after Recover still forks both branches.

### Edge Cases

- Single unconditional outgoing: unchanged (one token).
- Multiple outgoings with any condition or activity `default`: exclusive chooser (013), not inclusive fan-out.
- Parallel / exclusive / inclusive gateways: unchanged (they set Fork or OutgoingFlowID themselves).
- Join after activity fan-out uses existing parallel gateway join semantics.

## Requirements *(mandatory)*

- **FR-001**: Leaving an element with ≥2 outgoings, no conditions, and no default MUST take all outgoings (mint tokens for extras), emitting SEQUENCE_FLOW_TAKEN per flow.
- **FR-002**: Conditional / default selection MUST remain exclusive (one flow).
- **FR-003**: Start Event multi-unconditional outgoings MUST fan out on CreateInstance.
- **FR-004**: Recover MUST preserve waiting tokens so post-Recover Complete fans out identically.
- **FR-005**: Existing parallel gateway and conditional-flow suites MUST remain green.

## Success Criteria *(mandatory)*

- **SC-001**: Activity fan-out fixture completes both branches and the process.
- **SC-002**: Start Event fan-out fixture completes both branches.
- **SC-003**: Conditional activity tests still take one path.
- **SC-004**: Recover then Complete fans out both branches.
