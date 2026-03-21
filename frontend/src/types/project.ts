export interface ProjectPolishFormData {
  projectName: string;
  role: string;
  techStack: string;
  projectBackground: string;
  responsibility: string;
}

export interface FollowupQuestionItem {
  question: string;
  answer: string;
}

export interface ProjectPolishOutput {
  id: string;
  resumeDescription: string;
  interviewDescription: string;
  highlights: string[];
  difficulties: string[];
  followups: FollowupQuestionItem[];
}

export interface ProjectHistoryItem {
  id: string;
  projectName: string;
  createdAt: string;
  summary: string;
}