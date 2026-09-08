# Contracts: Nested Event Sub-Process

No new RPCs or proto fields.

## Deploy

- Event Sub-Process MAY be nested inside another Event Sub-Process.
- Existing Event Sub-Process shape rules still apply.

## Runtime

1. Entering an Event Sub-Process arms Event Sub-Processes whose parent is that Event Sub-Process.
2. Nested trigger/interrupt/complete uses parent scope = outer Event Sub-Process id.
3. Leaving the outer Event Sub-Process scope disarms nested arms.

## Recover

Nested arms while outer ESP is active survive Recover; nested trigger afterward matches continuous run.
