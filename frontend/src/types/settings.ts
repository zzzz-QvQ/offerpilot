export interface SettingsState {
  defaultInterviewMode: '综合模拟' | '专项训练' | '项目拷打';
  followUpIntensity: '低' | '中' | '高';
  scoringStrictness: '宽松' | '标准' | '严格';
  modelName: string;
  streamOutput: boolean;
}