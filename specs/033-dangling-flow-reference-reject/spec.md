# Feature Specification: Dangling Flow Reference Deploy Rejection

**Feature Branch**: `main` (developed on main)

**Created**: 2026-09-16

**Status**: Draft

**Input**: User description: "补 deploy 期拒绝：顺序流指向未建模/非可执行 flow element 时直接拒绝（关掉 choreography 静默通过这条路，与 032 同形）。"

## User Scenarios & Testing *(mandatory)*

BPMN XML may contain flow elements Sparrow does not model — choreography constructs (`choreographyTask`, `subChoreography`, `callChoreography`) and any other element outside the executable-process subset. The XML parser drops them, so today a definition that routes a token through such an element deploys successfully and only fails at `CreateInstance` with `UNSUPPORTED_ELEMENT: NOT_FOUND: element "…"`. The same holds for a declared `<outgoing>` / `<incoming>` that names a sequence flow which does not exist. Deploy MUST reject both instead, matching the implicit throw event decision in spec 032: unsupported work is refused at deploy, not discovered at runtime.

### User Story 1 - Sequence flow to an unmodeled element (Priority: P1)

A process routes `Start → choreographyTask → End`. Deploy fails with `UNSUPPORTED_ELEMENT` naming the sequence flow and the unresolved element id, and no deployment is stored.

### User Story 2 - Dangling reference in any scope (Priority: P1)

The same rejection applies when the dangling endpoint is inside an embedded SubProcess, a Transaction, an Ad-Hoc SubProcess, or a called process, because the check runs over the whole compiled element index.

### User Story 3 - Declared incoming/outgoing naming no sequence flow (Priority: P2)

An activity declares `<outgoing>Flow_missing</outgoing>` while no `sequenceFlow` with that id exists. Deploy fails with `UNSUPPORTED_ELEMENT` naming the element and the missing flow id, instead of the token failing when it tries to leave the element.

### User Story 4 - Existing definitions keep deploying (Priority: P1)

Every current fixture keeps deploying: elements that legitimately have no flows (event sub-processes, compensation handlers, boundary events referenced via `attachedToRef`) are unaffected, because only flow endpoints and declared flow references are checked.

## Requirements *(mandatory)*

- **FR-001**: Deploy MUST reject a `sequenceFlow` whose `sourceRef` or `targetRef` does not resolve to an indexed element, with `UNSUPPORTED_ELEMENT` naming the flow and the unresolved id.
- **FR-002**: Deploy MUST reject a declared `<incoming>` / `<outgoing>` reference that does not resolve to an indexed sequence flow, with `UNSUPPORTED_ELEMENT` naming the element and the unresolved id.
- **FR-003**: The check MUST cover every scope indexed by the deployment, including embedded SubProcess, Transaction, Ad-Hoc SubProcess, Event Sub-Process, and called processes.
- **FR-004**: Rejection MUST be deterministic: the same definition MUST always fail on the same reference, so the error text is stable across deploys.
- **FR-005**: No currently supported definition may start failing; elements without sequence flows remain valid.

## Success Criteria *(mandatory)*

- **SC-001**: A process with `Start → choreographyTask → End` is rejected at Deploy, not at CreateInstance.
- **SC-002**: A dangling endpoint inside a SubProcess is rejected the same way.
- **SC-003**: A declared outgoing naming a missing sequence flow is rejected.
- **SC-004**: The full existing test suite still passes.
