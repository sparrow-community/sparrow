# Research: Multiple Same-Kind Boundaries

## 1. Deploy uniqueness

**Decision**: Key uniqueness by (activity, kind, discriminator): timer→boundary id (always unique); message→message name; signal→signal name; compensate→empty (one); error→code (existing).

## 2. Arming model

**Decision**: Add `WaitingBoundary` message + `repeated waiting_boundaries` on `ActivityPayload`. `attachBoundary` fills the list for every timer/message/signal; still sets legacy singular fields to the first of each kind for GetInstance/compat.

## 3. Projection

**Decision**: `Token.BoundaryWaits []BoundaryWait`; apply from `waiting_boundaries` on ACTIVATED; clear on complete/terminate. Collectors prefer BoundaryWaits, fall back to legacy singular fields.

## 4. Interrupt siblings

**Decision**: `BoundaryEventHandler.OnComplete` interrupting path terminates all Timer/Message/Signal boundaries on the host except the firing one (slice helpers).

## 5. Protocol

**Decision**: Additive proto field; run `protocol/proto/build.sh`.
