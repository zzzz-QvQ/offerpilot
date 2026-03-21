import { http } from '@/utils/request';
import { STORAGE_KEYS } from '@/constants/storage';
import type {
  CreateInterviewSessionResponse,
  InterviewSessionDetail,
  SubmitAnswerParams,
  SubmitAnswerResponse,
} from '@/types/interview';

const buildBaseURL = () => {
  const baseURL = import.meta.env.VITE_API_BASE_URL ?? '';
  if (!baseURL) {
    return window.location.origin;
  }

  if (/^https?:\/\//.test(baseURL)) {
    return baseURL;
  }

  return new URL(baseURL, window.location.origin).toString();
};

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
  getStreamUrl(sessionId: string) {
    const token = localStorage.getItem(STORAGE_KEYS.token) ?? '';
    const url = new URL(`/api/interview/session/${sessionId}/stream`, buildBaseURL());

    if (token) {
      url.searchParams.set('token', token);
    }

    return url.toString();
  },
};