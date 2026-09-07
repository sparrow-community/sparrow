# Research: Instantiate Event-Based Gateway

## 1. Semantics: CreateInstance still creates the instance

**Decision**: Keep `CreateInstance` as the only way to mint a process instance. `instantiate=true` means the exclusive event-based gateway is the **process entry** (no none startEvent); CreateInstance enters that gateway and arms catch races. Do **not** auto-create instances from PublishMessage/PublishSignal/FireDue solely because an instantiate EBG is deployed (spec FR-010).

**Rationale**: Matches Sparrow’s COMMAND model and projection bootstrap; true “first event creates instance” needs a subscription registry outside this increment.

**Alternatives considered**:
- Message/signal subscription start without CreateInstance — larger surface, deferred.
- Require a dummy none startEvent feeding instantiate EBG — contradicts BPMN instantiate (no incoming) and FR-002.

## 2. Process entry resolution

**Decision**: Extend `StartEventID` / `Deployment.StartEventID()` so that:
1. If the process has startEvents → existing behavior (first startEvent with outgoing).
2. Else if exactly one process-level exclusive instantiate EBG with no incoming and ≥2 catch outs → return that gateway id.
3. Else → error (`no startEvent` / no valid entry).

Call sites (`CreateInstance`, `redriveCreateInstance`, Call Activity callee start) keep calling `StartEventID`; instantiate-only callees become startable the same way if needed.

**Rationale**: Minimal API churn; Enter(entryID) already dispatches by element type.

**Alternatives considered**:
- New `ProcessEntryID()` with dual call-site updates — clearer naming but more churn; optional later rename.
- Index `InstantiateGatewayID` on Deployment at compile — fine as helper inside StartEventID.

## 3. Deploy validation

**Decision**:
- Remove blanket `instantiate is not supported` reject.
- For `Instantiate==true`: require exclusive (empty or Exclusive; reject Parallel); no incoming sequence flows (check Gateway.Incoming and SequenceFlows targeting the gateway); ≥2 outs to intermediate catch events (existing out rules).
- Process may have zero startEvents iff there is exactly one valid process-level instantiate exclusive EBG.
- Reject process that has both ≥1 startEvent and any instantiate EBG.
- Instantiate EBG inside subprocess: reject or leave unsupported (prefer reject instantiate nested in subprocess — process entry only).

**Rationale**: Spec FR-001–FR-005; keep mid-process non-instantiate EBG unchanged.

**Alternatives considered**:
- Allow startEvent + instantiate EBG mid-process with instantiate=false only — instantiate mid-process is invalid BPMN; reject any Instantiate with incoming or nested.

## 4. Runtime handler

**Decision**: No change to `EventBasedGatewayHandler` expected. Instant lifecycle + Fork to catches; exclusive sibling cancel via existing `EventBasedSiblings` on catch complete.

**Rationale**: Mid-process exclusive EBG already implements the race; entry is just “Enter gateway first.”

## 5. Recover

**Decision**: `redriveCreateInstance` already calls `dep.StartEventID()` then `Enter`; once StartEventID returns the instantiate gateway, Recover mid-race is the same as mid-process EBG wait (replay projections + unfinished catch/timer work).

**Rationale**: Spec FR-009; constitution II.

## 6. Protocol / API

**Decision**: No proto changes. Deploy/CreateInstance/FireDue/PublishMessage contracts unchanged at the wire level; only deploy acceptance and entry element type for instantiate-only processes change.

## 7. Fixtures

**Decision**: `m11_instantiate_ebg_*.bpmn` — happy path (timer+message), invalid shapes (incoming, parallel instantiate, startEvent+instantiate, single catch).

**Rationale**: Align with m2/m3 EBG fixtures and prior increment naming.
