# Implementation Plan

**Branch**: `022-instantiate-receive-event-gateway` | **Spec**: [spec.md](./spec.md)

## Summary

Accept instantiate receive task and parallel / multiple instantiate event-based gateway entries; allow gateway targets to include receive tasks; CreateInstance arms every instantiate entry; exclusive wins cancel cross-entry waits.

## Research

- Parallel instantiate: reuse mid-process parallel (no sibling cancel).
- Multiple exclusive instantiate entries: CreateInstance enters each gateway/receive; `InstantiateAlternativePeers` expands TerminateWaitingAt.
- Receive as gateway target: validation + receive Complete uses EventBasedSiblings + peers.

## Files

`processing/deploy/deploy.go`, `processing/deploy/start.go`, `processing/engine.go`, `processing/handlers/receive_task.go` / `waiting_task.go` / `intermediate_catch_event.go`, tests + fixtures `m26_*`, AGENTS/README.
