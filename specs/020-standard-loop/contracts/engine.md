# Engine contract: Standard Loop

No proto schema change.

## Deploy

Accepts `standardLoopCharacteristics`. Rejects coexistence with `multiInstanceLoopCharacteristics`.

## Commands

`Complete` (and job Activate/Complete for job-backed types) each iteration. Loop exit takes outgoing once.

## Recover

Replay restores `StandardLoopIteration`; mid-loop Complete continues.
