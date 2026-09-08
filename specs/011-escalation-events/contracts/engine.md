# Contract: Escalation engine behavior

## Deploy

- Accept `escalationEventDefinition` on intermediate throw, end, boundary, ESP start.
- Resolve code from `<escalation escalationCode>` / name / id.
- Reject mixed event definitions on the same throw/catch element (existing extraCatch rules).

## Throw

1. Emit `ESCALATION_THROWN` with `escalation_code`.
2. Walk scope: armed escalation ESP → escalation boundary on current scope activity → terminate/bubble like error (for interrupting) → parent.
3. If none: return; caller continues outgoing or completes end.

## Catch

- Interrupting boundary / ESP: same terminate semantics as error counterparts.
- Non-interrupting boundary: spawn outgoing token; do not terminate host.
- Non-interrupting ESP: existing NI ESP path (consume thrower end token if end; leave parallel work).

## Recover

Replay ledger; no special redrive for escalation throw COMMAND (throws are EVENT from Enter). ESP handler wait Completes after Recover.
