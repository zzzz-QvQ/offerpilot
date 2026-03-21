export type InterviewQuestionStatus = 'pending' | 'in_progress' | 'completed';
export type KnowledgeHitLevel = 'high' | 'medium' | 'low';
export type InterviewAgentStage =
  | 'session_created'
  | 'generating_question'
  | 'retrieving_knowledge'
  | 'scoring'
  | 'generating_followup'
  | 'completed';

export type InterviewSSEEventType = 'delta' | 'state' | 'reference' | 'done';

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
  isStreaming?: boolean;
}

export interface InterviewStatusItem {
  label: string;
  value: string;
}

export interface InterviewScoreItem {
  label: string;
  score: number;
}

export interface KnowledgeReferenceItem {
  questionId: string;
  title: string;
  category: string;
  score: number;
  snippet: string;
}

export interface InterviewSessionDetail {
  sessionId: string;
  currentQuestionId: string;
  questions: InterviewQuestionItem[];
  messages: InterviewMessageItem[];
  statusItems: InterviewStatusItem[];
  scoreItems: InterviewScoreItem[];
  knowledgeHits: KnowledgeReferenceItem[];
  agentStage?: InterviewAgentStage;
}

export interface CreateInterviewSessionResponse {
  sessionId: string;
}

export interface SubmitAnswerParams {
  content: string;
}

export interface SubmitAnswerResponse {
  accepted: boolean;
}

export interface DeltaEventPayload {
  content: string;
}

export interface StateEventPayload {
  currentQuestionId?: string;
  questions?: InterviewQuestionItem[];
  statusItems?: InterviewStatusItem[];
  scoreItems?: InterviewScoreItem[];
  agentStage?: InterviewAgentStage;
}

export interface RawKnowledgeReferencePayloadItem {
  questionId?: string;
  title?: string;
  category?: string;
  score?: number;
  snippet?: string;
  name?: string;
  summary?: string;
  level?: KnowledgeHitLevel;
}

export interface ReferenceEventPayload {
  knowledgeHits?: RawKnowledgeReferencePayloadItem[];
}

export interface DoneEventPayload {
  currentQuestionId?: string;
  questions?: InterviewQuestionItem[];
  statusItems?: InterviewStatusItem[];
  scoreItems?: InterviewScoreItem[];
  knowledgeHits?: RawKnowledgeReferencePayloadItem[];
  agentStage?: InterviewAgentStage;
  message?: InterviewMessageItem;
}

export interface DeltaEvent {
  type: 'delta';
  payload: DeltaEventPayload;
}

export interface StateEvent {
  type: 'state';
  payload: StateEventPayload;
}

export interface ReferenceEvent {
  type: 'reference';
  payload: ReferenceEventPayload;
}

export interface DoneEvent {
  type: 'done';
  payload: DoneEventPayload;
}

export type InterviewSSEEvent = DeltaEvent | StateEvent | ReferenceEvent | DoneEvent;

export interface InterviewSSEEventMap {
  delta: DeltaEventPayload;
  state: StateEventPayload;
  reference: ReferenceEventPayload;
  done: DoneEventPayload;
}