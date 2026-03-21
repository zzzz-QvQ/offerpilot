import { http } from '@/utils/request';
import type { ReviewRadarItem, ReviewSuggestionItem, ReviewSummaryItem, ReviewTrendItem, ReviewWeaknessItem } from '@/types/review';

export interface ReviewOverview {
  overallScore: number;
  summaryItems: ReviewSummaryItem[];
  radarItems: ReviewRadarItem[];
  weaknessItems: ReviewWeaknessItem[];
  trendItems: ReviewTrendItem[];
  suggestionItems: ReviewSuggestionItem[];
}

export const reviewApi = {
  getOverview() {
    return http.get<ReviewOverview>('/reviews/overview');
  },
};