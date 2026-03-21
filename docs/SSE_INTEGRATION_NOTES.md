# SSE Integration Notes

## Findings

1. EventSource lifecycle is mostly correct: `useSSEStream.ts` already exposes `start()` and `stop()`, and closes the connection on `done`, `error`, and component unmount.
2. The SSE event chain is wired end to end:
   - `delta` updates the streaming interviewer message
   - `state` updates question/state/score/agent stage
   - `reference` updates knowledge hits
   - `done` finalizes the round
3. `InterviewStatusPanel` refreshes in real time because it consumes `agentStage`, `statusItems`, and `scoreItems` from the store, which are updated from `state`/`done` events.
4. `KnowledgeReferenceCard` refreshes in real time because it consumes `knowledgeHits`, which are updated from `reference`/`done` events.

## Issues Found

1. On submit failure, the optimistic candidate message remained in the local list, leaving frontend state ahead of backend state.
2. On manual stop or SSE error, `agentStage` could stay in an in-progress phase, making the status panel appear active after the stream had already ended.
3. When a stream was interrupted before any delta arrived, an empty interviewer placeholder bubble could remain in the chat list.
4. `useSSEStream.start()` returned a boolean in practice but did not declare that contract explicitly.

## Fixes Applied

1. `useSSEStream.ts`
   - Declared `start(url): boolean` explicitly.
   - Kept duplicate-connection protection for the same active URL.
   - Preserved cleanup on `done`, parse failure, connection failure, and unmount.

2. `stores/modules/interview.ts`
   - Added `rollbackPendingRound()` to remove the optimistic candidate message and streaming placeholder when submit fails.
   - Updated `interruptStreaming()` to remove an empty streaming placeholder and reset `agentStage` to `completed`.
   - Kept `streamError` and `streamStoppedByUser` as the UI-facing stream control state.

## Current Boundary

- REST is still responsible for:
  - creating a session
  - fetching the full session snapshot
  - accepting the candidate answer
  - finishing the session
- SSE is responsible for:
  - live interviewer text chunks
  - live agent stage updates
  - live score/status updates
  - live knowledge references
  - final round completion payload

## Residual Notes

1. The backend `POST /api/interview/session/:id/message` still returns a full payload shape even though the frontend now mainly uses it as an acknowledgment trigger before opening SSE. This is workable, but later it can be reduced to a lighter ack response if desired.
2. The current SSE handler does not send heartbeat events. The integration works, but heartbeat support would improve long-lived connection stability.