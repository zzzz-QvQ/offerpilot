import { http } from '@/utils/request';
import type { ReviewHistoryItem, ReviewReportResponse } from '@/types/review';

export const reviewApi = {
  getReport(sessionId: string) {
    return http.get<ReviewReportResponse>(`/review/${sessionId}`);
  },
  getHistory() {
    return http.get<ReviewHistoryItem[]>('/review/history');
  },
};