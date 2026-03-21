export type InterviewQuestionStatus = '待开始' | '进行中' | '已完成';
export type KnowledgeHitLevel = '高' | '中' | '低';

export interface InterviewQuestionItem {
  id: string;
  title: string;
  status: InterviewQuestionStatus;
}

export interface InterviewMessageItem {
  id: string;
  role: 'interviewer' | 'candidate';
  content: string;
  time: string;
}

export interface InterviewStatusItem {
  label: string;
  value: string;
}

export interface InterviewScoreItem {
  label: string;
  score: number;
}

export interface KnowledgeHitItem {
  name: string;
  level: KnowledgeHitLevel;
  summary: string;
}

export interface InterviewSessionDetail {
  sessionId: string;
  currentQuestionId: string;
  questions: InterviewQuestionItem[];
  messages: InterviewMessageItem[];
  statusItems: InterviewStatusItem[];
  scoreItems: InterviewScoreItem[];
  knowledgeHits: KnowledgeHitItem[];
}

export interface CreateInterviewSessionResponse {
  sessionId: string;
}

export interface SubmitAnswerParams {
  content: string;
}

export interface SubmitAnswerResponse {
  message: InterviewMessageItem;
  currentQuestionId: string;
  questions: InterviewQuestionItem[];
  statusItems: InterviewStatusItem[];
  scoreItems: InterviewScoreItem[];
  knowledgeHits: KnowledgeHitItem[];
}