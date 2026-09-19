# Plan

**On**: `main` | **Spec**: [spec.md](./spec.md)

Keep source BPMN on `deploy.Deployment` at Deploy and load; add `engine.v1.GetDeployment`; gateway adapter; rename `validateM1` → `validateProcess`. No ledger/subject changes.

Files: `processing/deploy/deploy.go`, `processing/engine.go`, `processing/open.go`, `protocol/proto/engine/v1/engine.proto`, `gateway/engine_server.go`, tests, AGENTS both editions.
