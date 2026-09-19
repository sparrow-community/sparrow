<!--
Sync Impact Report
- Version change: 1.0.1 → 1.0.2
- Modified principles: I wording (execution engine → execution and fact kernel) to match AGENTS; MUST intent unchanged
- Added sections: none
- Removed sections: none
- Follow-up TODOs: none
-->

# Sparrow Constitution

## Core Principles

### I. Executable Completeness, Not a Product Suite

Sparrow MUST remain a single-node BPMN **execution and fact kernel**. Coverage MUST grow toward full executable BPMN semantics; a permanent minimal subset is not the end state. The core MUST NOT grow a modeler, operations console, or other Camunda-style product-suite surfaces. Those MAY exist later as separate consumers of the engine, never as runtime dependencies.

Rationale: the project exists to execute and audit contracts, not to become a workflow product suite.

### II. Append-Only Ledger Is the Source of Truth

Every accepted process behavior MUST be recorded as an append-only Event. Projections MAY be discarded and MUST be rebuildable from the Event log plus deployments. COMMAND processing that crashes mid-chain MUST be finishable by Recover without inventing a new COMMAND. Job leases, late-message buffers, and similar helpers MUST live outside the ledger (optional runtime store). They MUST NOT become ledger subjects alongside Element.

Rationale: audit, replay, and recovery depend on one fact stream.

### III. Element Is the Subject of Behavior

An Event describes **what happened**. The nested Element (Type, id, token_id, Intent) is **who** it happened to. `Element.Type` MUST be one value per independent BPMN element (`PROCESS` plus concrete FlowElements). The engine MUST NOT collapse distinct BPMN elements into a generic TASK/GATEWAY plus a kind field. Job, timer, message, and signal waiting are payloads and lifecycle stages of that element, not peer subjects on the ledger.

Rationale: keeps the log aligned with BPMN and avoids a second ontology.

### IV. No Parallel Graph; Thin Module Boundaries

The runtime MUST execute the parsed `bpmn` definition (`element.Process`) after deploy-time validation. It MUST NOT maintain a second executable graph model. Transport (gRPC) MUST stay in `gateway`. Wire types MUST stay in `protocol` (Protocol Buffers only; generated Go is committed and not hand-edited). New executable elements MUST be added as a handler, registry entry, and deploy-time validation — not as Engine lifecycle branches.

Rationale: extra wrappers and crossed module duties have already been rejected as design debt.

### V. Standard BPMN; Generated Behavior Maps Onto It

AI-driven or custom behavior MUST map onto existing BPMN constructs (for example Service Task, User Task, and extensions). The engine MUST NOT invent non-standard core element types in the OMG sense. Definitions are the contract; the engine executes, validates, and rejects. It MUST NOT silently rewrite a definition because an agent proposed a change.

Rationale: generation (agent or otherwise) is variable; the contract is not silently mutable. See `AGENTS.md` Purpose and axioms.

### VI. Serial Commands, Explicit Rejection

COMMAND handling for a given `process_instance_id` MUST be strictly serial. Invalid or unsupported work MUST append a REJECTION with a stable code and message, and MUST NOT mutate instance projection as if the command succeeded. IDs for deployments, instances, and tokens MUST be UUIDv7 strings.

Rationale: one lock per instance is the partition model; silent success destroys auditability.

## Module Boundaries

The workspace is four Go modules. New top-level modules MUST NOT be added unless an existing module cannot host the work.

| Module | MUST | MUST NOT |
|--------|------|----------|
| `bpmn` | Parse and export BPMN 2.0 XML into `element` types | Encode runtime token or ledger semantics |
| `protocol` | Own `.proto` sources and generated packages (`event.v1`, `job.v1`, `engine.v1`) | Treat `job.v1` / `engine.v1` as Event record types |
| `processing` | Execute, recover, project; persist via `EventLog` + `deploy.Store` + optional `runtime.Store` | Host gRPC or other transports |
| `gateway` | Adapt EngineService and JobService; process entry `cmd/sparrow` | Reimplement element handlers |

License headers MUST remain Apache-2.0 (Sparrow community). Protocol changes MUST keep wire compatibility in mind (`buf` breaking category `FILE`).

## Quality Gates

- Semantic changes in `processing` MUST ship with tests that exercise the BPMN fixture and assert Event intents (not only the final projection).
- `go test ./processing/ ./gateway/ ./protocol/proto/event/v1/` MUST pass before a change is treated as done. The `bpmn` MIWG suite MAY run separately when XML model files change.
- Regenerating protocol MUST go through `protocol/proto/build.sh` (`buf lint` then `buf generate`). Hand-edited `*.pb.go` is forbidden.
- Planning for new work MUST proceed spec → plan → tasks against this constitution. `AGENTS.md` (English, canonical) and `AGENTS.zh.md` (Chinese) hold the implemented-element snapshot and roadmap index; they are not a second source of governing rules. Material snapshot changes MUST update both editions in the same change.

## Governance

This constitution supersedes informal notes when they conflict. Runtime design details live in `processing/README.md` and MUST remain consistent with these principles; if they diverge, amend this document or the design in the same change.

Amendments:

1. Record the change in this file (principle text, version, date).
2. Bump **Version** using semver: MAJOR for removed or redefined principles; MINOR for new or materially expanded principles; PATCH for clarification only.
3. Set **Last Amended** to the amendment date. Keep **Ratified** as the original adoption date.
4. Reviewers MUST check PRs against these principles (ledger subject, module boundary, no product-suite scope creep, tests for semantics).

Complexity that violates a principle (extra modules, parallel graph, new ledger subjects) MUST be listed in the feature plan's Complexity Tracking table with a simpler alternative that was rejected. Unjustified violations MUST block the plan.

**Version**: 1.0.2 | **Ratified**: 2026-08-26 | **Last Amended**: 2026-09-19
