# Feature Specification: Standalone Intermediate Escalation Catch

**Feature Branch**: `main` (developed on main)

**Created**: 2026-09-15

**Status**: Draft

**Input**: User description: "继续推进 Planned：Standalone intermediate escalation catch。保持在主分支开发。说明与文档中不要使用简写。"

## User Scenarios & Testing *(mandatory)*

Escalation throw, boundary, and event sub-process already work. Intermediate catch events still reject `escalationEventDefinition` at Deploy. This increment accepts a standalone intermediate escalation catch that waits until a matching escalation is thrown in its scope (or bubbles to that scope), then continues.

### User Story 1 - Parallel wait catches escalation throw (Priority: P1)

Fork: one branch waits on an intermediate escalation catch; the other completes work and throws a matching escalation. The catch completes, the join proceeds, and the process finishes.

**Acceptance Scenarios**:

1. **Given** parallel catch + throw with matching escalation code, **When** the throw fires, **Then** the waiting catch COMPLETED and the process can complete.
2. **Given** a catch with a non-matching code, **When** a different escalation is thrown in scope, **Then** the catch remains waiting (throw continues uncaught or caught elsewhere).

### User Story 2 - Catch-all escalation catch (Priority: P2)

An intermediate escalation catch with empty escalationRef (catch-all) receives any thrown escalation code in scope.

### User Story 3 - Recover while catch is waiting (Priority: P2)

Recover while the catch is waiting; after Recover the sibling throw still delivers and the process completes.

### Edge Cases

- Escalation event sub-process and escalation boundary in the same scope take precedence over intermediate catches (existing handlers first).
- Uncaught escalation still does not terminate the process.
- Event-based gateway may target an intermediate escalation catch (same target rules as other catches).
- Documentation spells out full terms without abbreviation.

## Requirements *(mandatory)*

- **FR-001**: Deploy MUST accept intermediateCatchEvent with exactly one escalationEventDefinition.
- **FR-002**: OnEnter MUST wait with EscalationCode in the ACTIVATED payload.
- **FR-003**: propagateEscalation MUST complete waiting intermediate escalation catches in the current scope whose code matches (exact or catch-all), after event sub-process and boundary checks for that scope.
- **FR-004**: Existing escalation suites MUST remain green.
- **FR-005**: Recover MUST restore a waiting intermediate escalation catch so a later throw still delivers.

## Success Criteria *(mandatory)*

- **SC-001**: Parallel fixture completes via catch after throw.
- **SC-002**: Code mismatch leaves catch waiting.
- **SC-003**: Recover mid-wait then throw → completed.
- **SC-004**: `m15_escalation_*` suites still pass.
