# Plan

**On**: `main` | **Spec**: [spec.md](./spec.md)

After indexing (process, embedded scopes, called processes) validate the compiled index: every `seqFlows` entry must resolve both endpoints, and every element's declared incoming/outgoing must name an indexed sequence flow. Iterate ids in sorted order so the first failure is stable. Run it from `compile` after `indexScope` / start indexing, so all scopes are present.

Files: `processing/deploy/flow_reference.go`, `processing/deploy/deploy.go`, fixtures `m37_*`, `processing/flow_reference_test.go`, AGENTS/README.
