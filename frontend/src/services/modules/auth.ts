import type { LoginParams, LoginResponse } from '@/types/user';

const MOCK_TOKEN_PREFIX = 'offerpilot_mock_token';

const mockLogin = async (params: LoginParams): Promise<LoginResponse> => {
  await new Promise((resolve) => {
    window.setTimeout(resolve, 600);
  });

  return {
    token: `${MOCK_TOKEN_PREFIX}_${Date.now()}`,
    user: {
      id: 'mock-user-1',
      email: params.email,
      nickname: 'OfferPilot User',
    },
  };
};

const loginByApi = async (params: LoginParams): Promise<LoginResponse> => {
  return mockLogin(params);
};

export const authApi = {
  login: loginByApi,
};