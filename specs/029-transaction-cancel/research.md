# Research

**Decision**: New `TYPE_TRANSACTION` (constitution: one Type per independent BPMN element).

**Decision**: Cancel = terminate unfinished in-scope tokens (keep host + cancel end) → compensation queue → fire interrupting cancel boundary.

**Decision**: Nest Transaction and MI Transaction rejected at Deploy; only `##Compensate`.
