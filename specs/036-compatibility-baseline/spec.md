# Feature Specification: Compatibility Baseline

**Feature Branch**: `036-compatibility-baseline`

**Created**: 2026-10-16

**Status**: Draft

**Input**: User description: "C-route locked; backward compatibility is core. Deliver P1 Compatibility baseline: release versioning scheme and 1.0 criteria; inventory of public surfaces; golden Event logs replayed in CI; buf breaking against last release (bootstrap if no tag yet)."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Know what is frozen (Priority: P1)

An embedder (or future maintainer) opens one document and sees: which surfaces are compatibility-bound, how alpha versions work today, what must be true before a 1.0 tag, and that REJECTION codes are part of the contract.

**Why this priority**: Without an explicit inventory, Principle VII cannot be reviewed in PRs.

**Independent Test**: The document lists Event wire, gRPC services, Go embed API, WASM host API, and REJECTION codes; README points at it and no longer says the wire is freely mutable without process.

**Acceptance Scenarios**:

1. **Given** the repository, **When** a reader opens the compatibility document, **Then** they can name every public surface and the release versioning stance.
2. **Given** a proposed removal of a REJECTION code, **When** a reviewer checks the inventory, **Then** they treat it as a breaking change under Principle VII.

### User Story 2 - Golden Event log survives Recover (Priority: P1)

A committed golden data directory (events.log + deployments) produced under the current kernel can be opened with `processing.Open` and rebuilds a known projection (status, waiting element, variables). CI runs this on every push.

**Why this priority**: Ledger replay is the core of Principle VII; CI must catch silent projection drift.

**Independent Test**: `go test` in `processing` fails if golden replay yields a different waiting element or status.

**Acceptance Scenarios**:

1. **Given** the golden baseline directory, **When** `Open` recovers it, **Then** the instance status and waiting user-task match the committed expectation.
2. **Given** a change that alters projection of that log without an upgrade path, **When** CI runs, **Then** the golden test fails.

### User Story 3 - Proto breaking checks in CI (Priority: P1)

Any PR that removes or renumbers a field in `event.v1` / `engine.v1` / `job.v1` fails CI via `buf breaking`. Until a published protocol tag exists, CI compares against a committed bootstrap snapshot of the current `.proto` tree.

**Why this priority**: Event wire is the hardest surface to repair after a break.

**Independent Test**: Temporarily deleting a field number fails `buf breaking` against the bootstrap; restoring it passes.

**Acceptance Scenarios**:

1. **Given** the bootstrap baseline matching current protos, **When** CI runs `buf breaking`, **Then** it passes on an unchanged tree.
2. **Given** a FIELD-category break, **When** CI runs, **Then** the job fails.

## Requirements *(mandatory)*

- **FR-001**: Repository MUST contain a compatibility document covering versioning stance (WASM calendar alpha today; Go modules untagged until an explicit release), 1.0 criteria, and a public-surface inventory (Event wire, `engine.v1`/`job.v1`, exported `processing` embed API, `globalThis.sparrow` + `sparrow.d.ts`, REJECTION codes).
- **FR-002**: Root README MUST stop claiming unrestricted wire mutability and MUST point at the compatibility document and constitution Principle VII.
- **FR-003**: At least one golden Event-log fixture MUST be committed and replayed in `processing` tests on every CI run of that module.
- **FR-004**: CI MUST run `buf breaking` for `protocol/proto` against a last-release baseline when a protocol tag exists; otherwise against a committed bootstrap snapshot of the proto tree at this increment.
- **FR-005**: AGENTS Planned P1 row MUST be updated to reflect delivery (or remaining bootstrap caveats).
- **FR-006**: A script or documented command MUST regenerate the golden fixture when maintainers intentionally refresh the baseline (not required on every CI run).

## Success Criteria *(mandatory)*

- **SC-001**: A new contributor can state the versioning stance and name the frozen surfaces from the compatibility document alone.
- **SC-002**: Golden replay test passes on mainline CI; a deliberate projection-breaking edit would fail it.
- **SC-003**: `buf breaking` runs in CI and fails on a FILE-category break against the bootstrap (or published) baseline.
- **SC-004**: Bootstrap path for “no published protocol tag yet” is documented so the next release can switch the against-target without redesign.

## Assumptions

- No Go module or protocol release tags exist yet in this clone; WASM continues calendar alpha (`wasm-vYYYY.M.D-alpha.N`).
- Golden fixture covers one representative in-flight instance (user-task wait), not every Supported combination—that remains the role of recover coverage tests.
- Product-suite / MCP adapters stay out of scope (consumers).
