export interface InterviewQuestionItem {
  id: string;
  title: string;
  status: '待开始' | '进行中' | '已完成';
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
  level: '高' | '中' | '低';
  summary: string;
}

export interface InterviewSessionSnapshot {
  sessionId: string;
  currentQuestionId: string;
  questions: InterviewQuestionItem[];
  messages: InterviewMessageItem[];
  statusItems: InterviewStatusItem[];
  scoreItems: InterviewScoreItem[];
  knowledgeHits: KnowledgeHitItem[];
}