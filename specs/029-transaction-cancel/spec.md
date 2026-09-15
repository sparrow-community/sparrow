# Feature Specification: Transaction SubProcess and Cancel

**Feature Branch**: `main` (developed on main)

**Created**: 2026-09-15

**Status**: Draft

**Input**: User description: "继续推进 Planned：Transaction SubProcess and cancel。TYPE_TRANSACTION；仅 ##Compensate；Cancel End 必须配 Cancel Boundary；嵌套 Transaction 与 MI Transaction Deploy 拒绝。"

## User Scenarios & Testing *(mandatory)*

Transaction is a specialized embedded SubProcess. Success completion matches SubProcess. Cancel End inside the Transaction terminates unfinished work, runs compensation for completed activities in reverse order, then leaves via an interrupting Cancel Boundary on the Transaction.

### User Story 1 - Successful Transaction (Priority: P1)

A Transaction with an inner User Task completes normally; the host takes its outgoing sequence flow. No cancel boundary fires.

### User Story 2 - Cancel with compensation (Priority: P1)

Inside a Transaction: Task_A completes with a compensation handler; Task_B waits; a path reaches Cancel End. Completing toward cancel terminates Task_B, runs Undo_A, then continues via Cancel Boundary to End_cancelled.

### User Story 3 - Cancel with no subscriptions (Priority: P2)

Cancel End with no completed compensatable activities still terminates the Transaction and leaves via Cancel Boundary.

### User Story 4 - Deploy rejections (Priority: P1)

Reject: cancel end outside Transaction; cancel boundary not on Transaction; cancel end without cancel boundary; method other than ##Compensate; nested Transaction; multi-instance on Transaction.

### User Story 5 - Recover (Priority: P2)

Recover while waiting on a compensation handler during cancel; Complete handler finishes cancel boundary path.

### Edge Cases

- Conditional / default paths inside Transaction unchanged.
- Error boundaries inside Transaction keep existing error semantics (out of cancel path).
- ##Image / ##Store unsupported at Deploy.

## Requirements *(mandatory)*

- **FR-001**: Deploy MUST accept `transaction` as `TYPE_TRANSACTION` with method empty or `##Compensate`.
- **FR-002**: Successful Transaction completion MUST behave like embedded SubProcess (host outgoing).
- **FR-003**: Cancel End MUST terminate unfinished tokens in the Transaction (keep host + cancel-end token), run in-scope compensation handlers reverse-order, then fire the Cancel Boundary and take its outgoing.
- **FR-004**: Cancel Boundary MUST attach only to a Transaction, MUST be interrupting, and MUST exist when any Cancel End is present in that Transaction.
- **FR-005**: Nested Transaction and multi-instance Transaction MUST be rejected at Deploy.
- **FR-006**: Recover MUST restore cancel-in-compensation state so remaining handlers complete into the cancel boundary path.

## Success Criteria *(mandatory)*

- **SC-001**: Success fixture completes via Transaction outgoing.
- **SC-002**: Cancel fixture runs Undo then Cancel Boundary path.
- **SC-003**: Deploy rejection fixtures return UNSUPPORTED_ELEMENT.
- **SC-004**: Recover mid-cancel compensation reaches the same completed outcome.
