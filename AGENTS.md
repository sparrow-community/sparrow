# Sparrow

Apache-2.0 Go workspace for a lightweight BPMN workflow engine.

Sparrow is a single-node execution engine in early development. The intended core is an append-only event log (Protocol Buffers) as the source of truth for process execution and audit. BPMN coverage grows from a small executable subset; the runtime is not a Camunda product clone.

Why this project exists in an AI-first world: see [`AI-Driven-BPMN.md`](./AI-Driven-BPMN.md).

## Workspace

Go **1.26.5** workspace (`go.work`) with three modules:

| Module | Path | Role | Maturity |
|--------|------|------|----------|
| `bpmn` | `./bpmn` | BPMN 2.0 XML model parse/export | Most complete |
| `protocol` | `./protocol` | Event protocol (Protobuf + Buf) | Schema skeleton |
| `processing` | `./processing` | Runtime / event handlers | Stub |

```text
sparrow/
├── go.work                 # committed
├── bpmn/                   # model layer
├── protocol/
│   ├── proto/              # .proto sources + Buf config
│   └── gen/go/             # generated Go (committed)
└── processing/             # engine stub
```

- Commit `go.work`; ignore `go.work.sum` (see root `.gitignore`).
- Module path prefix: `github.com/sparrow-community/sparrow/...`.
- `processing` uses local `replace` for `bpmn` and `protocol`.

## bpmn

- Entry: `bpmn.Bpmn` — load Definitions from bytes/file/stream; `FindElementByID`.
- Types: `bpmn/element` (large BPMN element set with XML tags).
- Schemas: `bpmn/schema` (BPMN 2.0 XSD).
- Tests: BPMN MIWG-oriented roundtrip/export suites under `bpmn/` and `bpmn/test`.

## protocol

Event messages only use **Protocol Buffers**. There is no FlatBuffers path.

- Sources: `protocol/proto/event/v1/*.proto`
- Generated Go: `protocol/gen/go/event/v1` (package `eventv1`)
- Import: `github.com/sparrow-community/sparrow/protocol/gen/go/event/v1`

Current `Event` shape (high level):

- `RecordType`: COMMAND / EVENT / REJECTION
- Instance fields: `deployment_id`, `process_instance_id`, `process_version` (ids are UUIDv7 strings)
- Causation: `source_record_id`; rejections carry `Rejection{code,message}`
- Nested `Element`: `Intent`, `Type`, `id`, `token_id` (token id is UUIDv7 string)
- `Element.Type` = one value per independent BPMN element (`PROCESS` + `FlowElements` concrete types); no generic TASK/GATEWAY + kind
- Payloads are data-only groups: process / event / activity / gateway / sequence-flow / data

### Regenerate

Requires Buf CLI (project targets current stable; e.g. 1.72.x).

```shell
cd protocol/proto && ./build.sh
```

`build.sh` runs `buf lint` then `buf generate`. Plugin is pinned in `buf.gen.yaml` (`buf.build/protocolbuffers/go:v1.36.11`). Managed mode sets `go_package_prefix` to `github.com/sparrow-community/sparrow/protocol/gen/go`.

`protocol/gen` is committed. Do not hand-edit `*.pb.go`; change `.proto` and regenerate. Generated output must match `./build.sh`.

## processing

- Runtime module: see `processing/README.md` and `processing/DESIGN.md`.
- M1 engine available: Deploy / CreateInstance / CompleteUserTask with in-memory or file event log.
- Durable open: `processing.Open(dataDir)` persists `events.log` + `deployments/*.bpmn` and replays projections.
- Layout: `deploy` (holds `element.Process`, no parallel graph), `handlers/` (one file per element type), `executor`, `projection`, `log`.
- IDs via UUIDv7 (`NextID` / `MustNextID`).

## Commands

From repo root (with `go.work`):

```shell
go test ./processing/ ./protocol/proto/event/v1/
go test ./bpmn/   # large MIWG suite; slower
cd protocol/proto && ./build.sh
```

## Conventions

- License headers: Apache-2.0 (Sparrow community).
- Prefer extending the existing three-module layout over adding new top-level modules without need.
- Protocol evolution: keep wire compatibility in mind (`buf` breaking category is `FILE`).
- AI / custom behavior, if added later, should map onto existing BPMN constructs (e.g. Service Task + extensions), not invent non-standard core element types in the OMG sense.
