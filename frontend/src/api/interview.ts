import { http } from '@/utils/request';
import type { InterviewMessageItem, InterviewQuestionItem, InterviewScoreItem, InterviewStatusItem, KnowledgeHitItem } from '@/types/interview';

export interface InterviewSessionDetail {
  sessionId: string;
  currentQuestionId: string;
  questions: InterviewQuestionItem[];
  messages: InterviewMessageItem[];
  statusItems: InterviewStatusItem[];
  scoreItems: InterviewScoreItem[];
  knowledgeHits: KnowledgeHitItem[];
}

export interface SubmitAnswerParams {
  sessionId: string;
  questionId: string;
  content: string;
}

export interface InterviewStreamParams {
  sessionId: string;
}

export const interviewApi = {
  getSession(sessionId: string) {
    return http.get<InterviewSessionDetail>(`/interviews/${sessionId}`);
  },
  createSession() {
    return http.post<{ sessionId: string }>('/interviews');
  },
  submitAnswer(data: SubmitAnswerParams) {
    return http.post<{ message: InterviewMessageItem; statusItems: InterviewStatusItem[]; scoreItems: InterviewScoreItem[] }>('/interviews/answer', data);
  },
  connectStream(_params: InterviewStreamParams) {
    return null;
  },
};