import { http } from '@/utils/request';
import type { ProjectHistoryItem, ProjectPolishFormData, ProjectPolishOutput } from '@/types/project';

export const projectApi = {
  generate(data: ProjectPolishFormData) {
    return http.post<ProjectPolishOutput>('/projects/polish', data);
  },
  getHistory() {
    return http.get<ProjectHistoryItem[]>('/projects/history');
  },
};