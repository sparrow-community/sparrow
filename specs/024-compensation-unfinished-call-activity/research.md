# Research

**Decision**: Run unfinished Call Activity inner compensation on the child instance; parent throw waits via WaitingChildren.

**Decision**: Reuse ResumeParent Completed=false to terminate the Call Activity host after child compensation (same as error/cancel paths).

**Decision**: Broadcast still applies unfinished SubProcess entry and same-scope collection; unfinished Call Activity is an additional entry kind.

**Decision**: Spell out “Call Activity” in AGENTS/README (no abbreviations).
