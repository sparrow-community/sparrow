# Plan

**On**: `main` | **Spec**: [spec.md](./spec.md)

Add `TYPE_TRANSACTION`, parse `transaction` / `cancelEventDefinition`, deploy rules, TransactionHandler (SubProcess-like), cancel path via PendingCompensation + CancelBoundary fields, fire cancel boundary after handlers.

Proto change + `build.sh`. Fixtures `m33_*`, tests, AGENTS/README.
