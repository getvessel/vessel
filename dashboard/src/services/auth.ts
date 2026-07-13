import type {
  AuthCredentials,
  AuthResponse,
  RegisterCredentials,
  Setup2FAResponse,
  Verify2FARequest,
} from '#/interfaces/auth';
import type { BaseResponse } from '#/interfaces/base';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const authService = {
  login: async (credentials: AuthCredentials): Promise<AuthResponse> => {
    try {
      return await apiClient.post<AuthResponse>('/auth/signin', credentials);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  register: async (details: RegisterCredentials): Promise<AuthResponse> => {
    try {
      return await apiClient.post<AuthResponse>('/auth/signup', details);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  logout: async (): Promise<void> => {
    try {
      await apiClient.post('/auth/logout');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  setup2FA: async (): Promise<BaseResponse<Setup2FAResponse>> => {
    try {
      return await apiClient.post<BaseResponse<Setup2FAResponse>>('/auth/2fa/setup');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  verify2FA: async (payload: Verify2FARequest): Promise<BaseResponse<void>> => {
    try {
      return await apiClient.post<BaseResponse<void>>('/auth/2fa/verify', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  disable2FA: async (payload: Verify2FARequest): Promise<BaseResponse<void>> => {
    try {
      return await apiClient.post<BaseResponse<void>>('/auth/2fa/disable', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
