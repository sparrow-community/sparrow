# Quickstart

1. Deploy parallel caller + called process (Task_A compensatable, Task_B waits).
2. CreateInstance → complete Task_A on child → leave Task_B → complete sibling → compensate throw.
3. Complete Undo_A on child → caller completes.
4. Recover mid Undo_A and Complete → same final status.
