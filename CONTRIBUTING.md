# Contributing to Sparrow

Thanks for helping with Sparrow. This repository is the BPMN **execution and fact kernel**. Agents, UIs, and other tools are consumers of the command and event contract—they are not part of this tree.

## Public repositories and sync

| Role | Location |
|------|----------|
| Public source of truth for readers | [github.com/sparrow-community/sparrow](https://github.com/sparrow-community/sparrow) |
| Primary maintainer workspace | Cursor Origin remotes under the `slowrookie` account (this checkout’s `origin`) |

Maintainers develop in Cursor with `origin` as the working remote. GitHub
(`sparrow-community/sparrow`) is the public mirror for readers, Actions, and npm
Trusted Publishing. Do not treat a personal fork of Origin as the public project home.

**Automatic mirror:** `.github/workflows/mirror-from-origin.yml` (schedule +
manual) fetches Origin `main`/tags and pushes to GitHub. Requires Actions secrets
`ORIGIN_GIT_TOKEN` (read Origin) and `GH_MIRROR_TOKEN` (GitHub PAT with `repo` +
`workflow` so mirrored tags can trigger `publish-wasm`). Until secrets exist the
workflow no-ops with a warning.

**WASM npm release:** push tag `wasm-vYYYY.M.D-alpha.N` on Origin; after mirror,
`publish-wasm` builds and publishes. See [`wasm/README.md`](./wasm/README.md).

Sibling playground: [sparrow-playground](https://github.com/sparrow-community/sparrow-playground).

## Before you change semantics

Read [`AGENTS.md`](./AGENTS.md) (Chinese: [`AGENTS.zh.md`](./AGENTS.zh.md)) and [`.specify/memory/constitution.md`](./.specify/memory/constitution.md). Runtime notes live in [`processing/README.md`](./processing/README.md).

- Effective definition and instance state change only through accepted **COMMAND**s; every accepted behavior is an append-only **EVENT**.
- Combinations with runtime token or side-effect semantics are either **Supported** or **Excluded**. Unsupported work is rejected at Deploy. Silence is not exclusion.
- Do not invent ledger subjects, empty intents, or non-OMG core element types.

Feature work for new Supported combinations follows `/speckit-specify` → plan → tasks → implement (`.cursor/skills/`).

## Development setup

- Go **1.26.5**
- Full checkout including [`go.work`](./go.work) (there is no root `go.mod`)

```shell
go test -C bpmn ./...
go test -C processing ./...
go test -C gateway ./...
go test -C protocol ./...
go test -C wasm ./...
```

Edit `.proto` sources under `protocol/proto`, then regenerate with `cd protocol/proto && ./build.sh`. Do not hand-edit `*.pb.go`.

Gateway (trusted network only; no TLS or auth):

```shell
go run ./gateway/cmd/sparrow -data-dir ./data -listen :50051
```

WASM artifact for the playground:

```shell
cd wasm && ./build.sh
```

## Pull requests

1. Keep each PR focused on one logical change.
2. Include or update end-to-end and Recover tests when you change Supported element semantics.
3. Keep Apache-2.0 file headers on new Go sources.
4. Expect CI (`.github/workflows/ci.yml`) to run module-scoped `go test` on GitHub.

## License

By contributing, you agree that your contributions are licensed under the [Apache License 2.0](./LICENSE). See [`NOTICE`](./NOTICE) for third-party attribution (BPMN MIWG fixtures under CC BY 3.0).

## Conduct

Please follow the [Code of Conduct](./CODE_OF_CONDUCT.md).
