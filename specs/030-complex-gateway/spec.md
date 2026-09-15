# Feature Specification: Complex Gateway

**Feature Branch**: `main` (developed on main)

**Created**: 2026-09-15

**Status**: Draft

**Input**: User description: "继续推进 Planned：Complex Gateway。"

## User Scenarios & Testing *(mandatory)*

Complex Gateway supports merge with optional `activationCondition` and divergence with sequence-flow conditions (inclusive-style take-many). Empty activationCondition on a multi-incoming join waits for all incoming tokens (parallel-join semantics).

### User Story 1 - Activation join (Priority: P1)

Three parallel branches join at a Complex Gateway with `activationCondition` `${arrived >= 2}`. When two tokens have arrived, the gateway fires, terminates the waiting peer, and continues; the third branch may arrive later and start a new wait cycle.

### User Story 2 - Parallel-like join without condition (Priority: P1)

Complex Gateway with multiple incomings and no activationCondition waits for all incoming tokens, then continues.

### User Story 3 - Complex split (Priority: P1)

Complex Gateway with multiple outgoings evaluates conditions like Inclusive Gateway (all matching + unconditional; else default).

### User Story 4 - Recover (Priority: P2)

Recover while one token waits at a complex join; a second arrival after Recover still activates when the condition is met.

## Requirements *(mandatory)*

- **FR-001**: Deploy MUST accept `complexGateway` as `TYPE_COMPLEX_GATEWAY` with optional `activationCondition` and `default`.
- **FR-002**: Multi-incoming join with activationCondition MUST fire when the expression is true given `arrived` / `nrOfInstances` (arrived token count including current) and `incoming` (incoming flow count).
- **FR-003**: Multi-incoming join without activationCondition MUST wait for all incoming tokens.
- **FR-004**: On fire, waiting peer tokens at the gateway MUST be terminated; outgoing paths follow inclusive-style selection when multiple outgoings exist.
- **FR-005**: Multi-outgoing split (single incoming) MUST select outgoings like Inclusive Gateway.
- **FR-006**: Recover MUST preserve waiting tokens at the complex join.

## Success Criteria *(mandatory)*

- **SC-001**: `arrived >= 2` join fires after two arrivals.
- **SC-002**: No-condition join waits for all branches.
- **SC-003**: Split takes matching condition paths.
- **SC-004**: Recover then second arrival activates correctly.
