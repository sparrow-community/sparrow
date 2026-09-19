# Feature Specification: GetDeployment

**Feature Branch**: `main` (developed on main)

**Created**: 2026-09-19

**Status**: Draft

**Input**: User description: "继续推进 GetDeployment；顺便处理 deploy.go 中的 validateM1。"

## User Scenarios & Testing *(mandatory)*

Consumers overlay Events onto a deployment diagram. They need the original BPMN XML for a deployment id without keeping a private copy. Go already exposes compiled `GetDeployment`; gRPC does not return XML. Separately, `validateM1` is a leftover milestone name for process-level deploy validation and should be renamed so the kernel no longer refers to M1.

### User Story 1 - Fetch definition after Deploy (Priority: P1)

Deploy a process, then `GetDeployment` with that id returns the same BPMN bytes, plus `process_id` and `process_version`.

### User Story 2 - Fetch after Recover (Priority: P1)

`Open` / `Recover` reloads deployments from the deploy store; `GetDeployment` still returns the persisted XML for each id.

### User Story 3 - Missing id (Priority: P2)

Unknown deployment id returns `NOT_FOUND` (gRPC `NotFound`).

### User Story 4 - validateProcess rename (Priority: P2)

Deploy validation still rejects invalid boundaries and start/instantiate combinations; the entry point is no longer named `validateM1`.

## Requirements *(mandatory)*

- **FR-001**: After a successful Deploy, the engine MUST retain the submitted BPMN XML with the deployment.
- **FR-002**: `engine.v1.GetDeployment` MUST return `deployment_id`, `process_id`, `process_version`, and `bpmn_xml` for a known id.
- **FR-003**: After Recover/Open with a deploy store, GetDeployment MUST return the stored XML.
- **FR-004**: Unknown deployment id MUST fail with a stable `NOT_FOUND` (mapped to gRPC NotFound).
- **FR-005**: Process-level deploy validation formerly named `validateM1` MUST be renamed to a non-milestone name (`validateProcess`) without changing rejection semantics.

## Success Criteria *(mandatory)*

- **SC-001**: Deploy then GetDeployment round-trips XML bytes.
- **SC-002**: Open/Recover then GetDeployment returns persisted XML.
- **SC-003**: Unknown id is NotFound / NOT_FOUND.
- **SC-004**: Existing deploy rejection tests still pass after the rename.
