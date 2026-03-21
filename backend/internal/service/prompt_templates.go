package service

import "fmt"

func buildInterviewInitialQuestionSystemPrompt() string {
	return "You are a senior frontend interviewer. Generate one concise but realistic opening interview question. Return plain text only."
}

func buildInterviewInitialQuestionUserPrompt(mode, jobRole string) string {
	return fmt.Sprintf(
		"Interview mode: %s\nJob role: %s\nGenerate the first interview question for this session.",
		mode,
		jobRole,
	)
}

func buildInterviewEvaluationSystemPrompt() string {
	return "You are an interview evaluator. Return strict JSON only with keys summary, followupFocus, expressionScore, accuracyScore, structureScore."
}

func buildInterviewEvaluationUserPrompt(sessionMode, jobRole, currentQuestion, answer, ragContext string) string {
	return fmt.Sprintf(
		"Job role: %s\nMode: %s\nCurrent question: %s\nCandidate answer: %s\nRetrieved knowledge:\n%s\nEvaluate the answer and return JSON only.",
		jobRole,
		sessionMode,
		currentQuestion,
		answer,
		ragContext,
	)
}

func buildInterviewFollowupSystemPrompt() string {
	return "You are a senior frontend interviewer. Write concise interview feedback followed by exactly one follow-up question. Plain text only."
}

func buildInterviewFollowupUserPrompt(sessionMode, jobRole, currentQuestion, answer, ragContext string, evaluation *interviewEvaluationResult) string {
	return fmt.Sprintf(
		"Job role: %s\nMode: %s\nCurrent question: %s\nCandidate answer: %s\nEvaluation summary: %s\nFollow-up focus: %s\nRetrieved knowledge:\n%s\nGenerate short interviewer feedback and then one follow-up question.",
		jobRole,
		sessionMode,
		currentQuestion,
		answer,
		evaluation.Summary,
		evaluation.FollowupFocus,
		ragContext,
	)
}

func buildProjectPolishSystemPrompt() string {
	return "You are a senior frontend interview coach and resume optimizer. Rewrite project experience into strong, factual job-hunting content. Return strict JSON only with keys: resumeVersion, interviewVersion, highlights, difficulties, followups. followups must be an array of objects with question and answer. Do not wrap in markdown."
}

func buildProjectPolishUserPrompt(input projectRawInput) string {
	return fmt.Sprintf(
		"Project Name: %s\nBackground: %s\nTech Stack: %s\nResponsibilities: %s\nDifficulties: %s\nAchievements: %s\nRole: %s\n\nRequirements:\n1. resumeVersion: concise, resume-ready project description in one paragraph.\n2. interviewVersion: interview-ready explanation with challenge and result.\n3. highlights: 3 to 5 bullet-style strings focusing on technical impact.\n4. difficulties: 2 to 3 strings describing difficulty and solution.\n5. followups: 3 likely interviewer questions with strong suggested answers.\nReturn JSON only.",
		input.ProjectName,
		input.Background,
		input.TechStack,
		input.Responsibilities,
		input.Difficulties,
		input.Achievements,
		input.Role,
	)
}

func buildReviewAnalysisSystemPrompt() string {
	return "You are a senior frontend interview reviewer. Read the interview session and return strict JSON only with keys: overallSummary, overallScore, technical, expression, logic, depth, project, weakPoints, suggestions. All scores must be integers between 0 and 100. weakPoints and suggestions must each contain 2 to 4 concise strings."
}

func buildReviewAnalysisUserPrompt(sessionMode, jobRole, conversation, referenceContext string) string {
	return fmt.Sprintf(
		"Session mode: %s\nJob role: %s\nConversation:\n%s\n\nRetrieved references:\n%s\n\nGenerate a structured review report in JSON.",
		sessionMode,
		jobRole,
		conversation,
		referenceContext,
	)
}
