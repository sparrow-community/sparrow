# Plan

**On**: `main` | **Spec**: [spec.md](./spec.md)

Accept CatchKindConditional on event sub-process starts; store condition text; arm on scope entry; EvaluateConditions triggers matching armed event sub-processes via existing `triggerEventSubProcess`.

Files: `deploy/event_subprocess.go`, `deploy/deploy.go`, `event_subprocess.go`, `conditional_eval.go`, fixtures `m30_*`, tests, AGENTS/README.
