# Contract: Link engine behavior

## Deploy

- Accept link intermediate throw/catch.
- Reject throw without ≥1 same-name catch.
- Reject throw with outgoing / catch with incoming sequence flows.

## Runtime

1. Enter throw → lifecycle + optional link_name payload.
2. Enter each matching catch (reuse token then mint).
3. Catch → lifecycle + TakeOutgoing.

## Recover

No special COMMAND; waiting after jump is normal User Task Complete.
