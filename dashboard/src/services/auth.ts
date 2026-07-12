import type { AuthCredentials, AuthResponse, RegisterCredentials } from '#/interfaces/auth';
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
};
