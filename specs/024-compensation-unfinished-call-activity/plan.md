# Implementation Plan

**Branch**: `024-compensation-unfinished-call-activity` | **Spec**: [spec.md](./spec.md)

## Summary

Extend `startCompensation` so parent-scope compensate throws enter unfinished Call Activity children: publish child compensation work, run handlers on the child instance, terminate child + Call Activity host, then advance parent PendingCompensation.

## Research

- Child subscriptions and handlers live on the child instance/deployment — cannot Enter them on the caller.
- New publication `compensate_unfinished_child`; flush runs child cancel + local compensation queue.
- Child PendingCompensation without a throw token finishes by terminating the instance and ResumeParent(Completed=false).
- Parent PendingCompensation.WaitingChildren gates advance until all unfinished-call children report done.

## Files

`processing/handlers/handler.go`, `executor.go`, `call_child.go` / `signals.go`, `projection/instance.go`, fixtures `m28_*`, tests, AGENTS/README.
