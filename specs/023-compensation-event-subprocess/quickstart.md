# Quickstart

1. Deploy BPMN with embedded SubProcess containing a compensation event sub-process (compensate start).
2. CreateInstance → complete SubProcess work → SubProcess completes (compensation subscribed).
3. Reach compensate throw → compensation event sub-process runs.
4. Complete waiting work inside it → throw continues → process completes.
5. Recover mid-wait and Complete → same final status.
