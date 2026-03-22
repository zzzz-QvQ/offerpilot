# Code Cleanup Notes

## Scope

This cleanup pass focused on removing code that is clearly unused now, without touching active API, SSE, LLM, or RAG boundaries.

## Removed

### 1. Unused frontend app store

Deleted:

- `frontend/src/stores/modules/app.ts`

Reason:

- The store was not imported or referenced anywhere in the project.
- It only contained sidebar collapse state that the current layout does not use.

### 2. Unused global placeholder style

Updated:

- `frontend/src/styles/index.scss`

Removed:

- `.placeholder-card`

Reason:

- The class was not referenced by any current view or component.

### 3. Unused interview store method

Updated:

- `frontend/src/stores/modules/interview.ts`

Removed:

- `setError()`

Reason:

- The method had no call sites.
- Error handling already uses direct store state updates and `setStreamError()`.

## Checked But Intentionally Kept

### 1. Interview SSE compatibility normalization

Kept in:

- `frontend/src/stores/modules/interview.ts`
- `frontend/src/types/interview.ts`

Reason:

- The `RawKnowledgeReferencePayloadItem` compatibility layer still protects the frontend from backend protocol drift.
- This is still useful boundary code, not dead code.

### 2. Backend placeholder-style fallbacks inside real services

Kept in:

- `backend/internal/service/dashboard_service.go`
- `backend/internal/service/question_service.go`

Reason:

- These are not mock data modules anymore.
- They are runtime-safe fallbacks when upstream data is missing, which is acceptable for a deliverable system.

### 3. Prompt template and service boundaries

Kept in:

- `backend/internal/service/prompt_templates.go`
- `backend/internal/service/llm_service.go`
- `backend/internal/service/question_retrieval_service.go`

Reason:

- These are part of the active production structure and should remain.

## Validation

Passed after cleanup:

- `frontend`: `npm.cmd run build`
- `backend`: `go build ./...`
