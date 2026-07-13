import type {
  ListEnabledProvidersResponse,
  ListOAuthProvidersResponse,
  SaveOAuthProviderRequest,
  SaveOAuthProviderResponse,
} from '#/interfaces/oauth';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const oauthService = {
  listProviders: async (): Promise<ListOAuthProvidersResponse> => {
    try {
      const response = await apiClient.get<ListOAuthProvidersResponse>(
        '/api/v1/settings/oauth/providers'
      );
      return response;
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getEnabledProviders: async (): Promise<ListEnabledProvidersResponse> => {
    try {
      const response = await apiClient.get<ListEnabledProvidersResponse>(
        '/api/v1/auth/oauth/providers/enabled'
      );
      return response;
    } catch (error) {
      throw handleApiError(error);
    }
  },

  saveProvider: async (payload: SaveOAuthProviderRequest): Promise<SaveOAuthProviderResponse> => {
    try {
      const response = await apiClient.put<SaveOAuthProviderResponse>(
        '/api/v1/settings/oauth/providers',
        payload
      );
      return response;
    } catch (error) {
      throw handleApiError(error);
    }
  },

  triggerOAuthLogin: (provider: string) => {
    window.location.href = `/api/v1/auth/oauth/${provider}`;
  },
};
