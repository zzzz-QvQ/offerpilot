import { http } from '@/utils/request';
import type { PageParams, PageResult } from '@/types/api';
import type { QuestionBankFilters, QuestionItem } from '@/types/question-bank';

export interface QuestionListParams extends Partial<QuestionBankFilters>, Partial<PageParams> {}

export const questionApi = {
  getList(params?: QuestionListParams) {
    return http.get<PageResult<QuestionItem>>('/questions', { params });
  },
  getDetail(questionId: string) {
    return http.get<QuestionItem>(`/questions/${questionId}`);
  },
  toggleFavorite(questionId: string) {
    return http.post<void>(`/questions/${questionId}/favorite`);
  },
};