# Compatibility

Sparrow’s core promise (constitution **Principle VII**): facts outlive releases. An Event log written by a published release must Recover on every later release; public surfaces evolve additively.

This document is the inventory and release stance for that promise. Element coverage lives in [`AGENTS.md`](./AGENTS.md). Runtime design lives in [`processing/README.md`](./processing/README.md).

## Versioning stance (today)

| Artifact | Scheme | Notes |
|----------|--------|-------|
| WASM npm `@sparrow-community/wasm` | Calendar alpha `YYYY.M.D-alpha.N` via tags `wasm-vYYYY.M.D-alpha.N` | Both `latest` and `alpha` dist-tags move on every publish |
| Go modules (`bpmn`, `protocol`, `processing`, `gateway`, `wasm`) | Untagged / pseudo-versions until an explicit module tag | No `v1.0.0` yet |
| Gateway binary | Same commit as the modules it embeds | Not separately versioned |
| Protocol wire (`event.v1`, `engine.v1`, `job.v1`) | Additive under `buf` FILE breaking; CI against bootstrap until a protocol tag exists | See [Proto breaking](#proto-breaking) |

**Alpha means:** behavior covered by tests is real; the surfaces below are already treated as compatibility-bound for review. A break still requires an upgrade path and a constitution-reviewed justification—even before 1.0.

## 1.0 criteria

A `v1.0.0` (Go modules and a matching WASM non-alpha line) MAY be cut only when all of the following hold:

1. Golden Event logs from every prior published release (or the committed bootstrap set) replay cleanly on the candidate.
2. `buf breaking` against the previous published protocol tag is green.
3. Public surfaces listed below have had no unannounced breaks for at least one release cycle.
4. README and this file no longer describe the wire as a free-for-all preview.
5. AGENTS Supported set remains covered by end-to-end and Recover tests.

Until then, stay on calendar alpha for WASM; do not imply SemVer 1.0 stability on Go module tags.

## Public surfaces (compatibility-bound)

### Event wire (`event.v1`)

- Package: `protocol/proto/event/v1/`
- On-disk log: length-delimited protobuf records (`processing/log` File backend)
- Record types: COMMAND, EVENT, REJECTION
- Breaking: remove/renumber/retype/repurpose fields or enum values; change delimited framing

### gRPC (`engine.v1`, `job.v1`)

- `EngineService` / `JobService` in `protocol/proto/engine/v1/`, `job/v1/`
- Breaking: remove RPCs or request/response fields; rename REJECTION codes observed by clients

### Go embed API (`processing`)

Compatibility-bound entry points for embedders:

- Lifecycle: `NewEngine`, `Recover`, `Open`, `(*Engine).Close`
- Commands: `Deploy`, `CreateInstance`, `CreateInstanceRequest`, `Complete`, `ThrowError`, `ResolveIncident`, `Activate`, `Fail`, `Heartbeat`
- Async / typed start: `FireDue`, `NextDueUnixMs`, `SetNow`, `PublishMessage`, `PublishSignal`, `EvaluateConditions`, `EvaluateConditionalStarts`
- Query: `GetDeployment`, `GetDeploymentXML`, `GetInstance`, `ListInstanceIDs`, `ListEvents`
- Intervention: `EnableIntervention`, `DisableIntervention`, `SetBreakpoints`, `Continue`, `StepInto`, `StepOver`, `Pause`, `HostEffectAllowed`, `GetInterventionState`, `SetVariables`
- Persistence contracts: `log.EventLog`, `deploy.Store`, `runtime.Store` and their file/memory constructors

Unexported helpers and test-only symbols are not a contract.

### WASM host (`globalThis.sparrow` + `sparrow.d.ts`)

Methods and types in [`wasm/sparrow.d.ts`](./wasm/sparrow.d.ts) (deploy through intervention APIs). Errors use `{ $error: string }`. There is no engine API version field today; npm package version is the release label.

### REJECTION / error codes

Stable string codes. Do not rename or reuse.

| Code | Typical use |
|------|-------------|
| `INVALID_ARGUMENT` | Malformed command input |
| `INVALID_CONDITION` | Timer/duration/condition evaluation |
| `INVALID_EXPRESSION` | Expression parse/eval |
| `INVALID_MAPPING` | Call Activity IO mapping |
| `INVALID_STATE` | Wrong lifecycle / not waiting |
| `NOT_FOUND` | Missing instance, element, deployment, version |
| `NO_OUTGOING_FLOW` | No eligible outgoing sequence flow |
| `UNSUPPORTED_ELEMENT` | Deploy/runtime unsupported executable construct |
| `INCIDENT_OPEN` | Work blocked on open incident |
| `NO_INCIDENT` | Resolve/retry when no incident |
| `INTERNAL` | Fallback when an error lacks a structured code |

Ledger REJECTION records carry `code` + `message` on `Event_RECORD_TYPE_REJECTION`.

## Golden Event logs

Committed under `processing/testdata/compat/`. Each baseline directory is a full `Open` data dir (`events.log`, `deployments/`, optional `runtime/`) plus `expected.json`.

CI runs `TestCompatGoldenBaseline` in the `processing` module. To refresh intentionally after an additive, reviewed change:

```shell
./processing/scripts/regen-compat-baseline.sh
```

Refreshing because a test failed without an upgrade path is a Principle VII violation—fix the change or provide a migrator.

## Proto breaking

- Config: `protocol/proto/buf.yaml` → `breaking.use: [FILE]`
- **Bootstrap (no protocol release tag yet):** committed snapshot at `protocol/breaking-baseline/`. CI runs:

  ```shell
  cd protocol/proto && buf breaking --against ../breaking-baseline
  ```

- **After a protocol tag exists:** switch CI to  
  `buf breaking --against '.git#tag=<tag>,subdir=protocol/proto'`  
  and keep or retire the bootstrap directory in that same change.

Local check (requires [buf](https://buf.build)):

```shell
cd protocol/proto && buf breaking --against ../breaking-baseline
```

## Related

- Constitution Principle VII: [`.specify/memory/constitution.md`](./.specify/memory/constitution.md)
- Spec: [`specs/036-compatibility-baseline/`](./specs/036-compatibility-baseline/)
