import { http } from '@/utils/request';
import type {
  CreateInterviewSessionResponse,
  InterviewSessionDetail,
  SubmitAnswerParams,
  SubmitAnswerResponse,
} from '@/types/interview';

export const interviewApi = {
  createSession() {
    return http.post<CreateInterviewSessionResponse>('/interview/session');
  },
  getSession(sessionId: string) {
    return http.get<InterviewSessionDetail>(`/interview/session/${sessionId}`);
  },
  submitMessage(sessionId: string, data: SubmitAnswerParams) {
    return http.post<SubmitAnswerResponse>(`/interview/session/${sessionId}/message`, data);
  },
  finishSession(sessionId: string) {
    return http.post<void>(`/interview/session/${sessionId}/finish`);
  },
};