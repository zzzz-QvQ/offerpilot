export interface UserProfile {
  id: string;
  email: string;
  nickname: string;
}

export interface LoginParams {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: UserProfile;
}