# Plan

**On**: `main` | **Spec**: [spec.md](./spec.md)

Tests only, reusing existing fixtures (`m2_inclusive_gateway.bpmn`, `m25_mi_none_behavior_event.bpmn` + `m25_catch_mi_each.bpmn`, `m26_instantiate_receive.bpmn`, `m5_version_v1/v2.bpmn`, `m4_error_boundary.bpmn`). Each test drives the engine with `processing.Open(dir)`, closes it, reopens the same data dir, and continues the scenario. Any defect the tests reveal is fixed in the runtime under the same spec.

Files: `processing/recover_coverage_test.go`, AGENTS/README if a runtime fix is needed.
