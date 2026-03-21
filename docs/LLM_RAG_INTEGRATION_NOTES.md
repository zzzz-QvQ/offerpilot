# LLM / RAG Integration Notes

## Scope

This note summarizes the current fifth-phase integration status for:

- Unified LLM service
- Embedding + Milvus retrieval chain
- Interview real LLM + RAG loop
- Project Polish real generation
- Review real analysis

## Current Findings

1. The unified `LLMService` is already reused by `Interview`, `Project`, and `Review`.
2. The retrieval chain is structurally clear:
   - `EmbeddingService` handles `/embeddings`
   - `QuestionRetrievalService` handles question vectorization and search
   - `MilvusStore` handles collection creation, upsert, and search
3. `Interview` has already moved from placeholder generation to a real loop:
   - create session
   - generate first question with LLM
   - submit candidate answer
   - retrieve question knowledge
   - evaluate answer with LLM
   - stream feedback and follow-up through SSE
4. `Project` uses the unified LLM service for real content generation.
5. `Review` uses real session messages plus optional retrieval context to build analysis through the unified LLM service.

## Fixes Made In This Pass

### 1. Interview reference payload aligned with frontend expectations

The backend `reference` event now uses the same practical shape the frontend prefers:

- `questionId`
- `title`
- `category`
- `score`
- `snippet`

This reduces normalization ambiguity between backend SSE payloads and frontend interview store handling.

### 2. Retrieval failures now degrade gracefully

When embedding or Milvus is not configured, `Interview` and `Review` no longer fail the whole request flow immediately.

Current fallback behavior:

- `Interview`: returns empty reference hits and continues the LLM path
- `Review`: skips reference context and continues review generation

This keeps the business loop usable while the vector layer is being prepared.

### 3. Review service cleanup

The review service file was cleaned up to remove damaged string literals and inconsistent suggestion type labels.

Current suggestion item types are now:

- `suggestion`
- `weakness`

## Current Module Boundaries

### LLM

- File: `backend/internal/service/llm_service.go`
- Responsibilities:
  - standard text generation
  - streaming text generation
  - OpenAI-compatible request handling

### Retrieval

- File: `backend/internal/service/embedding_service.go`
- File: `backend/internal/service/question_retrieval_service.go`
- File: `backend/internal/pkg/vectorstore/milvus.go`
- Responsibilities:
  - embedding generation
  - vector collection management
  - question indexing and semantic retrieval

### Interview

- File: `backend/internal/service/interview_service.go`
- Responsibilities:
  - session lifecycle
  - first-question generation
  - answer evaluation
  - retrieval-backed prompt context
  - SSE event emission

### Project

- File: `backend/internal/service/project_service.go`
- Responsibilities:
  - normalize raw input
  - call LLM for structured generation
  - persist generated project polish result

### Review

- File: `backend/internal/service/review_service.go`
- Responsibilities:
  - load persisted report when available
  - otherwise derive report from interview session + messages
  - optionally enrich prompt context with retrieval results

## SSE Contract Status

The Interview SSE contract remains:

- `state`
- `reference`
- `delta`
- `done`

The main alignment point after this pass is that `reference` now better matches frontend display needs.

## Remaining Gaps

1. `Project` and `Review` still use `context.Background()` internally instead of request-scoped context.
2. `Review` trend data is still synthetic rather than fully derived from persisted historical reports.
3. Retrieval fallback is now safe, but collection bootstrap and data freshness still depend on separate indexing flow.
4. Structured LLM output parsing is still duplicated as loose JSON parsing conventions rather than a more formal shared schema helper.

## Recommended Next Steps

1. Pass request-scoped context into `ProjectService` and `ReviewService` generation paths.
2. Persist more structured evaluation outputs from `Interview` for stronger downstream review history.
3. Add a small shared helper for structured LLM JSON output validation.
4. Add operational notes for Milvus collection initialization and question re-indexing workflow.
