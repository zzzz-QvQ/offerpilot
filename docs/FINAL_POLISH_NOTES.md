# Final Polish Notes

## Problem List

Ordered by severity:

1. Multiple user-facing pages and components contained mojibake text and a few broken template fragments.
2. Some transient UI states could survive route changes:
   - Question drawer visibility
   - Question detail loading state
   - Interview stream warning / stream error hint
3. A few stores were missing complete error handling for detail fetch or favorite actions.
4. Settings and question-bank type definitions still contained corrupted literal values, which made maintenance harder.
5. Several fallback error messages were inconsistent or unreadable.

## Fixes Made

### User-facing text and template cleanup

Cleaned and normalized visible copy in:

- `frontend/src/views/login/LoginView.vue`
- `frontend/src/views/settings/SettingsView.vue`
- `frontend/src/components/project-polish/ProjectInputForm.vue`
- `frontend/src/components/question-bank/QuestionFilterBar.vue`
- `frontend/src/components/question-bank/QuestionDetailDrawer.vue`
- `frontend/src/components/question-bank/QuestionList.vue`

### Type cleanup

Normalized literal types in:

- `frontend/src/types/question-bank.ts`
- `frontend/src/types/settings.ts`

### Store and composable polish

Improved state handling in:

- `frontend/src/stores/modules/dashboard.ts`
- `frontend/src/stores/modules/question-bank.ts`
- `frontend/src/stores/modules/review.ts`
- `frontend/src/stores/modules/project.ts`
- `frontend/src/stores/modules/settings.ts`
- `frontend/src/composables/useQuestionBank.ts`
- `frontend/src/composables/useInterviewSession.ts`
- `frontend/src/composables/useProjectPolish.ts`

Key changes:

- Added clearer fallback error messages
- Added question-bank transient state reset on unmount
- Cleared interview stream feedback on mount/unmount
- Added error handling for question detail and favorite toggling
- Added history-fetch failure handling in project store

## Current Status

### Loading

Core loading states are present across the main modules:

- Login submit loading
- Dashboard page loading
- Question list / detail / favorite loading
- Interview session / submit / finish / stream loading
- Review page loading
- Project generate / history / detail loading

### Error

Error states are present in the main stores and views. The most visible unreadable fallback strings have been corrected.

### Empty

Main empty states are already present for:

- Dashboard
- Question list
- Question detail
- Review
- Project result/history-dependent areas

### Duplicate request / duplicate submit / duplicate SSE

Current protections already in place:

- Login prevents repeated submit while loading
- Interview blocks submit during stream
- SSE composable prevents reopening the same active stream URL
- Interview stop logic closes the current stream before finish

### State residue

Improved in this pass:

- Question-bank transient drawer/loading/error state no longer lingers after page switch
- Interview stream feedback is cleared on mount/unmount

## Remaining Non-blocking Notes

1. Question search still triggers fetches directly from filter changes. This is acceptable now, but can be optimized later with debounce.
2. Frontend build still reports large chunks because of ECharts and current bundling strategy.
3. Some modules still intentionally keep store state across navigation for continuity, which is acceptable and not treated as a bug.

## Verification

- Frontend build passed with `npm.cmd run build`
