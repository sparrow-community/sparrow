# Plan

**On**: `main` | **Spec**: [spec.md](./spec.md)

Model `AdHocSubProcess` in bpmn; add `TYPE_AD_HOC_SUB_PROCESS` to the event protocol; index inner activities plus `completionCondition` / `ordering` / `cancelRemainingInstances` at deploy with a flat-body rule; `AdHocSubProcessHandler` parks a host token and the executor enables inner activities, re-evaluates the completion condition on each inner completion (cancel remaining, enable next, or exhaust), then completes the scope on the host token.

Files: `bpmn/element/ad_hoc_sub_process.go`, `protocol/proto/event/v1/event.proto`, `processing/deploy/ad_hoc_sub_process.go`, `processing/handlers/ad_hoc_sub_process.go`, `processing/ad_hoc_executor.go`, fixtures `m35_*`, tests, AGENTS/README.
