# Research: Escalation Events

## 1. Semantics vs error

**Decision**: Same scope walk (ESP → scope boundary → parent → …). Uncaught escalation does **not** call terminate-instance; intermediate throw continues outgoing; escalation end runs TryCompleteProcess.

**Rationale**: BPMN / Camunda: uncaught escalation is ignored; uncaught error fails the process.

## 2. Ledger

**Decision**: Add `INTENT_ESCALATION_THROWN` and `EventPayload.escalation_code`. Do not overload `ERROR_THROWN`.

## 3. Boundary interrupting

**Decision**: Honor `cancelActivity` (explicit true/false in fixtures). Non-interrupting: mint token for boundary outgoing; leave host/thrower as appropriate.

## 4. Root element

**Decision**: Add `element.Escalation` with `EscalationCode`; register on `RootElemnts`. Resolve like errors.

## 5. Activity ThrowEscalation API

**Decision**: Out of scope; throws only from intermediate throw / end (handler `ThrowEscalation` effect).
