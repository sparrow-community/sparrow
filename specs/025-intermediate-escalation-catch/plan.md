# Plan

**On**: `main` | **Spec**: [spec.md](./spec.md)

Accept escalation intermediate catch at Deploy; wait on Enter; complete matching waits from `propagateEscalation` after event sub-process / boundary checks in the same scope.

Files: `deploy/escalation.go`, `deploy/deploy.go`, `deploy/timer.go` (CatchKind), `handlers/intermediate_catch_event.go`, `escalations.go`, fixtures `m29_*`, tests, AGENTS/README.
