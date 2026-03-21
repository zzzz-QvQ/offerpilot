export interface ReviewSummaryItem {
  label: string;
  value: string;
  hint: string;
}

export interface ReviewRadarItem {
  name: string;
  score: number;
  fullMark: number;
}

export interface ReviewWeaknessItem {
  name: string;
  score: number;
}

export interface ReviewTrendItem {
  date: string;
  score: number;
}

export interface ReviewSuggestionItem {
  title: string;
  description: string;
  type: '建议' | '薄弱点';
}