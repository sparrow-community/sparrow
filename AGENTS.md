# Sparrow

Apache-2.0 Go workspace for a lightweight BPMN workflow engine.

Sparrow is a single-node execution engine in early development. The intended core is an append-only event log (Protocol Buffers) as the source of truth for process execution and audit. BPMN coverage grows from a small executable subset; the runtime is not a Camunda product clone.

Why this project exists in an AI-first world: see [`AI-Driven-BPMN.md`](./AI-Driven-BPMN.md).

## Workspace

Go **1.26.5** workspace (`go.work`) with four modules:

| Module | Path | Role | Maturity |
|--------|------|------|----------|
| `bpmn` | `./bpmn` | BPMN 2.0 XML model parse/export | Most complete |
| `protocol` | `./protocol` | Protobuf: ledger + job/engine RPC | event / job / engine v1 |
| `processing` | `./processing` | Runtime / event handlers | M1 + Job pull |
| `gateway` | `./gateway` | gRPC adapter + `cmd/sparrow` | JobService + EngineService |

```text
sparrow/
├── go.work                 # committed
├── bpmn/                   # model layer
├── protocol/
│   ├── proto/              # .proto sources + Buf config
│   └── gen/go/             # generated Go (committed)
├── processing/             # engine
└── gateway/                # gRPC + cmd/sparrow
```

- Commit `go.work`; ignore `go.work.sum` (see root `.gitignore`).
- Module path prefix: `github.com/sparrow-community/sparrow/...`.
- `processing` uses local `replace` for `bpmn` and `protocol`.
- `gateway` uses local `replace` for `processing`, `protocol`, and `bpmn`.

## bpmn

- Entry: `bpmn.Bpmn` — load Definitions from bytes/file/stream; `FindElementByID`.
- Types: `bpmn/element` (large BPMN element set with XML tags).
- Schemas: `bpmn/schema` (BPMN 2.0 XSD).
- Tests: BPMN MIWG-oriented roundtrip/export suites under `bpmn/` and `bpmn/test`.

## protocol

Event messages only use **Protocol Buffers**. There is no FlatBuffers path.

- Sources: `protocol/proto/event/v1/*.proto`, `protocol/proto/job/v1/*.proto`, `protocol/proto/engine/v1/*.proto`
- Generated Go: `event/v1` (`eventv1`), `job/v1` (`jobv1`), `engine/v1` (`enginev1`) under `protocol/gen/go`

`event.v1` is the append-only ledger. `job.v1` (`JobService`) is the worker command RPC. `engine.v1` (`EngineService`) is the process-client RPC (Deploy / CreateInstance / Complete / query). Neither RPC package is an Event record type.

Current `Event` shape (high level):

- `RecordType`: COMMAND / EVENT / REJECTION
- Instance fields: `deployment_id`, `process_instance_id`, `process_version` (ids are UUIDv7 strings)
- Causation: `source_record_id`; rejections carry `Rejection{code,message}`
- Nested `Element`: `Intent`, `Type`, `id`, `token_id` (token id is UUIDv7 string)
- `Element.Type` = one value per independent BPMN element (`PROCESS` + `FlowElements` concrete types); no generic TASK/GATEWAY + kind
- Payloads are data-only groups: process / event / activity / gateway / sequence-flow / data
- `EventPayload` may carry `due_unix_ms` + original timer text (`duration` field, timeDuration / timeDate / timeCycle) for timer catch ACTIVATED, `message_name` for message catch ACTIVATED, and `variables` on catch COMPLETING

### Regenerate

Requires Buf CLI (project targets current stable; e.g. 1.72.x).

```shell
cd protocol/proto && ./build.sh
```

`build.sh` runs `buf lint` then `buf generate`. Plugins are pinned in `buf.gen.yaml` (`buf.build/protocolbuffers/go:v1.36.11`, `buf.build/grpc/go:v1.6.2`). Managed mode sets `go_package_prefix` to `github.com/sparrow-community/sparrow/protocol/gen/go`.

`protocol/gen` is committed. Do not hand-edit `*.pb.go`; change `.proto` and regenerate. Generated output must match `./build.sh`.

## processing

- Runtime module: see `processing/README.md` and `processing/DESIGN.md`.
- M1 engine available: Deploy / CreateInstance / Complete with in-memory or file event log.
- Waiting activities (UserTask, ServiceTask, intermediate timer catch, intermediate message catch) complete via one `Complete` API; type comes from the deployment. ServiceTask job type is `ActivityPayload.job_type`. Timer due is `EventPayload.due_unix_ms` from `timeDuration`, `timeDate`, or `timeCycle` (first due only); `FireDue` completes expired timer catches. Message catch name is `EventPayload.message_name`; `PublishMessage` completes matching waiters by name plus optional `correlation_keys` (instance variable JSON match; no buffer). `cmd/sparrow` ticks `FireDue`.
- Workers pull ServiceTask jobs with `Activate` (in-memory lease; not an EventLog record). `Fail` keeps the token waiting and releases the lease; `Heartbeat` extends it. Complete still finishes the waiting token.
- External clients use `gateway` (`engine.v1` + `job.v1` gRPC). Process entry is `gateway/cmd/sparrow` (`Open` + Serve). Do not put transport in `processing`.
- Durable open: `processing.Recover(ctx, eventLog, deploymentStore)` rebuilds projections and finishes any COMMAND whose EVENT chain was interrupted. `Open(ctx, dataDir)` is the file-backed convenience.
- Persistence is pluggable: `log.EventLog` + `deploy.Store` (Memory/File in-tree; swap in your own).
- Layout: `deploy` (holds `element.Process`, no parallel graph), `handlers/` (one file per element type), `executor`, `projection`, `log`.
- IDs via UUIDv7 (`NextID` / `MustNextID`).

## Commands

From repo root (with `go.work`):

```shell
go test ./processing/ ./gateway/ ./protocol/proto/event/v1/
go test ./bpmn/ # large MIWG suite; slower
cd protocol/proto && ./build.sh
go run ./gateway/cmd/sparrow -data-dir ./data -listen :50051
```

## Conventions

- License headers: Apache-2.0 (Sparrow community).
- Prefer extending the existing module layout over adding new top-level modules without need (`gateway` is the transport adapter; do not put gRPC in `processing`).
- Protocol evolution: keep wire compatibility in mind (`buf` breaking category is `FILE`).
- AI / custom behavior, if added later, should map onto existing BPMN constructs (e.g. Service Task + extensions), not invent non-standard core element types in the OMG sense.
