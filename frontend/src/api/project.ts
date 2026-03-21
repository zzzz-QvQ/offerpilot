import { http } from '@/utils/request';
import type { ProjectHistoryItem, ProjectPolishFormData, ProjectPolishOutput } from '@/types/project';

export const projectApi = {
  generate(data: ProjectPolishFormData) {
    return http.post<ProjectPolishOutput>('/project/polish', data);
  },
  getHistory() {
    return http.get<ProjectHistoryItem[]>('/project/polish/history');
  },
  getDetail(id: string) {
    return http.get<ProjectPolishOutput>(`/project/polish/${id}`);
  },
};