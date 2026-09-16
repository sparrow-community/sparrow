# Feature Specification: Recover Coverage for Remaining Supported Combinations

**Feature Branch**: `main` (developed on main)

**Created**: 2026-09-16

**Status**: Draft

**Input**: User description: "做一轮 Supported 组合的语义完整性审计…补齐四项 Recover 测试，并额外补上「重启后新发 ThrowError」那条 error boundary 路径。"

## User Scenarios & Testing *(mandatory)*

AGENTS states completeness holds only when every Supported combination has end-to-end **and** Recover tests. An audit of `processing` (test functions that restart the engine via `processing.Open` or `processing.Recover`) found end-to-end coverage complete, but five Recover paths missing or weak. This increment closes them with tests only; no runtime behavior is expected to change. A test that fails reveals a real Recover defect and is fixed in the runtime.

### User Story 1 - Inclusive Gateway across restart (Priority: P1)

Two branches are taken; one completes so a token waits at the inclusive join. After restart the waiting join token is rebuilt and the second branch completing still fires the join and completes the instance.

### User Story 2 - Multi-instance behavior events across restart (Priority: P1)

A multi-instance activity with a `noneBehaviorEventRef` signal publishes on each inner completion. After restart, completing a further inner still publishes to a catcher created after the restart, and the loop still completes.

### User Story 3 - Instantiate receive task entry across restart (Priority: P1)

An instance entered through an instantiate receive task waits for its message. After restart the message subscription is rebuilt and `PublishMessage` still completes the instance.

### User Story 4 - Coexisting revisions across restart (Priority: P1)

Revision 1 and revision 2 of the same process id both have in-flight instances. After restart each instance stays bound to its own revision and completes on its own path; `ProcessID` alone still resolves to the latest revision and an explicit `ProcessVersion` still resolves to revision 1.

### User Story 5 - Error boundary with a fresh throw after restart (Priority: P2)

An activity with an interrupting error boundary is waiting. After restart a **new** `ThrowError` command (not a redriven one) is still caught by the rearmed boundary. This complements the existing redrive test, which replays a partially applied `ERROR_THROWN`.

## Requirements *(mandatory)*

- **FR-001**: A Recover test MUST cover a token waiting at an inclusive join, asserting the join still fires after restart.
- **FR-002**: A Recover test MUST cover multi-instance behavior event publishing after restart.
- **FR-003**: A Recover test MUST cover an instantiate receive task entry, asserting `PublishMessage` still completes the instance after restart.
- **FR-004**: A Recover test MUST cover coexisting revisions, asserting per-instance revision binding, latest resolution, and explicit version resolution after restart.
- **FR-005**: A Recover test MUST cover a fresh `ThrowError` against a rearmed interrupting error boundary after restart.
- **FR-006**: Tests MUST assert Event intents where the behavior is an intent (boundary fired, signal delivered), not only the final projection status.

## Success Criteria *(mandatory)*

- **SC-001**: All five tests pass against the current runtime, or a revealed Recover defect is fixed in the runtime and the test then passes.
- **SC-002**: Every Supported combination in AGENTS has at least one end-to-end and one Recover test after this increment.
