# Engine contract: Abstract Task

No proto schema change (`TYPE_TASK` exists).

## Deploy

Accepts `task`. Indexes `TYPE_TASK`. May host boundaries; multi-instance like Manual Task.

## Commands

`Complete` on waiting abstract task (same as Manual). Not claimed by Job `Activate`.

## Recover

Same wait rebuild as Manual Task.
