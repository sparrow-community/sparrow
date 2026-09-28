# Sparrow

Sparrow is a single-node [BPMN 2.0](https://www.omg.org/spec/BPMN/2.0/) execution and fact kernel. It accepts a versioned process definition, runs that definition, and appends what happened to an event log. Effective definition and instance state change only through accepted commands. Projections of instance state can be rebuilt by replaying the log.

This repository is the kernel: BPMN XML model, protobuf contract, execution engine, gRPC gateway, and an in-memory WebAssembly host. Agents, diagram overlays, and operations tools are consumers of that contract. They are not part of this tree.

**Early public preview (alpha).** Behavior covered by the test suite is real, and the wire contract is still moving. Protobuf package names such as `engine.v1`, `event.v1`, and `job.v1` are not a compatibility freeze. RPC shapes, event payloads, and module paths can still change before a stable release.

Public repository: [github.com/sparrow-community/sparrow](https://github.com/sparrow-community/sparrow). Browser demo: [sparrow-playground](https://github.com/sparrow-community/sparrow-playground). How to contribute: [`CONTRIBUTING.md`](./CONTRIBUTING.md).

The supported executable-process subset, and the combinations this kernel refuses, are the agent contract in [`AGENTS.md`](./AGENTS.md) (Chinese: [`AGENTS.zh.md`](./AGENTS.zh.md)). Runtime notes live in [`processing/README.md`](./processing/README.md).

## Requirements

- Go **1.26.5**
- A checkout that includes [`go.work`](./go.work)

There is no root `go.mod`. The workspace lists five modules:

| Directory | Module path | Role |
|-----------|-------------|------|
| `bpmn` | `github.com/sparrow-community/sparrow/bpmn` | BPMN XML model |
| `protocol` | `github.com/sparrow-community/sparrow/protocol` | Protobuf sources and generated Go |
| `processing` | `github.com/sparrow-community/sparrow/processing` | Execution and fact kernel |
| `gateway` | `github.com/sparrow-community/sparrow/gateway` | gRPC server |
| `wasm` | `github.com/sparrow-community/sparrow/wasm` | In-memory browser host |

Sibling modules are wired with `replace` directives. Build and test from a full checkout of this repository.

## Tests

`go test ./...` from the repository root fails. That directory is not a module, so the workspace rejects the pattern:

```text
pattern ./...: directory prefix . does not contain modules listed in go.work or their selected dependencies
```

Run tests inside each module. From the repository root:

```shell
go test -C bpmn ./...
go test -C processing ./...
go test -C gateway ./...
go test -C protocol ./...
go test -C wasm ./...
```

The same commands work after `cd` into the module (`go test ./...`). A path such as `go test ./processing/` from the root tests only that one package and skips nested packages (`processing/expr`, `processing/log`, and the rest). `./...` has to be evaluated inside the module.

`protocol` tests use the generated code already committed under `protocol/gen/go`. Edit `.proto` sources only when changing the contract, then regenerate with `cd protocol/proto && ./build.sh`. Do not hand-edit `*.pb.go`.

## Run the gateway

From the repository root:

```shell
go run ./gateway/cmd/sparrow -data-dir ./data -listen :50051
```

`-data-dir` stores the event log and deployments (default `data/`). `-listen` is the gRPC address (default `:50051`). The process also fires due timers.

The listener has **no authentication and no TLS**. `gateway.NewServer` does not install transport credentials or auth interceptors. Bind it on a trusted network only. See [`SECURITY.md`](./SECURITY.md).

## WebAssembly demo

The browser demo is a separate repository, [sparrow-playground](https://github.com/sparrow-community/sparrow-playground). It is not published from this tree. Playground loads the artifact produced by the `wasm` module and schedules timers and jobs in JavaScript. The engine build itself stays in memory and does not open a port.

```shell
cd wasm && ./build.sh
# writes dist/sparrow.wasm (gitignored) and dist/wasm_exec.js
```

The JS host contract is [`wasm/README.md`](./wasm/README.md).

## License

Source code is [Apache License 2.0](./LICENSE). Copyright The Sparrow community and contributors.

The BPMN XML and PNG files under [`bpmn/test/`](./bpmn/test/) are unmodified BPMN MIWG fixtures (bpmn.io / Camunda Modeler results). They are Creative Commons Attribution 3.0 Unported, not Apache-2.0. See [`NOTICE`](./NOTICE) and [`bpmn/README.md`](./bpmn/README.md).
