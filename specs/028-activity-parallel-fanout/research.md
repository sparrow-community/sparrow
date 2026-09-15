# Research

**Decision**: Implicit AND-split only when every outgoing lacks a condition and the element has no `default` (BPMN parallel from activity / start). Mixed conditions stay exclusive per 013.

**Decision**: Reuse executor Fork machinery (token[0] continues, mint for rest) rather than inventing a new effect.
