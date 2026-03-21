export type QuestionFrequency = '高频' | '中频';

export interface QuestionCategoryOption {
  label: string;
  value: string;
}

export interface QuestionItem {
  id: string;
  title: string;
  category: string;
  frequency: QuestionFrequency;
  isFavorite: boolean;
  summary: string;
  difficulty: '基础' | '进阶';
  answerPoints: string[];
}

export interface QuestionBankFilters {
  keyword: string;
  category: string;
  frequency: '' | QuestionFrequency;
}