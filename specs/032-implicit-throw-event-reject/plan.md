# Plan

**On**: `main` | **Spec**: [spec.md](./spec.md)

Reject `implicitThrowEvent` FlowElements in `Deployment.indexScope`, which every scope (process, embedded SubProcess, Transaction, Ad-Hoc SubProcess, called process) passes through. Multi-instance behavior events live under `multiInstanceLoopCharacteristics` and are untouched. Record the element as Excluded in AGENTS with the rejection rationale.

Files: `processing/deploy/deploy.go`, fixtures `m36_*`, `processing/implicit_throw_event_test.go`, AGENTS/README.
