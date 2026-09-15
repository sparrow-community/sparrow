# Plan

**On**: `main` | **Spec**: [spec.md](./spec.md)

Add `Deployment.ChooseOutgoingFlows`. When TakeOutgoing has no preferred flow and multiple unconditional outs exist, reuse Enter Fork token minting from Complete / Enter / SubProcess leave.

Files: `processing/deploy/deploy.go`, `processing/executor.go`, fixtures `m32_*`, tests, AGENTS/README.
