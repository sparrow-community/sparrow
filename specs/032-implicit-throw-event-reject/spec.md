# Feature Specification: Implicit Throw Event Deploy Rejection

**Feature Branch**: `main` (developed on main)

**Created**: 2026-09-16

**Status**: Draft

**Input**: User description: "继续完成 Planned：Implicit throw event — reject at deploy with a stable code。"

## User Scenarios & Testing *(mandatory)*

`implicitThrowEvent` is in the BPMN `flowElement` substitution group, but as a process FlowElement it carries no token semantics: BPMN uses it inside `complexBehaviorDefinition` for multi-instance behavior events. Today Deploy parses such an element and then ignores it, so the definition deploys and a token entering it fails late at runtime. Deploy MUST reject it with a stable code instead, and it MUST NOT be aliased to a none throw event.

### User Story 1 - Process-level implicit throw (Priority: P1)

A process contains `implicitThrowEvent` between two tasks. Deploy fails with `UNSUPPORTED_ELEMENT` naming the element; no deployment is stored and no instance can be created.

### User Story 2 - Implicit throw inside any scope (Priority: P1)

The same rejection applies inside an embedded SubProcess, a Transaction, an Ad-Hoc SubProcess, an Event Sub-Process, and a called process, because every scope is indexed the same way.

### User Story 3 - Multi-instance behavior events keep working (Priority: P1)

`complexBehaviorDefinition` (and `noneBehaviorEventRef` / `oneBehaviorEventRef`) carry an implicit throw event inside `multiInstanceLoopCharacteristics`, not in `flowElements`. Those deployments MUST stay accepted and keep publishing their signal / message.

## Requirements *(mandatory)*

- **FR-001**: Deploy MUST reject a `implicitThrowEvent` that appears as a FlowElement of any process scope with an `UNSUPPORTED_ELEMENT` error naming the element id.
- **FR-002**: Deploy MUST NOT register such an element in the element index, and MUST NOT treat it as a none throw event or any other supported throw.
- **FR-003**: Multi-instance behavior events built from `implicitThrowEvent` inside `multiInstanceLoopCharacteristics` MUST remain unaffected.

## Success Criteria *(mandatory)*

- **SC-001**: Deploying a process with a process-level `implicitThrowEvent` returns `UNSUPPORTED_ELEMENT`.
- **SC-002**: The same rejection happens for an implicit throw nested in a SubProcess.
- **SC-003**: The existing multi-instance complex behavior fixtures still deploy and publish.
