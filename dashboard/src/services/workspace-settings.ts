import type {
  GetAISettingsResponse,
  GetEmailSettingsResponse,
  SaveAISettingsRequest,
  SaveAISettingsResponse,
  SaveEmailSettingsRequest,
  SaveEmailSettingsResponse,
} from '#/interfaces/workspace-settings';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const workspaceSettingsService = {
  getAISettings: async (workspaceId: string): Promise<GetAISettingsResponse> => {
    try {
      return await apiClient.get<GetAISettingsResponse>(`/workspaces/${workspaceId}/ai_settings`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  saveAISettings: async (
    workspaceId: string,
    payload: SaveAISettingsRequest
  ): Promise<SaveAISettingsResponse> => {
    try {
      return await apiClient.put<SaveAISettingsResponse>(
        `/workspaces/${workspaceId}/ai_settings`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getEmailSettings: async (workspaceId: string): Promise<GetEmailSettingsResponse> => {
    try {
      return await apiClient.get<GetEmailSettingsResponse>(
        `/workspaces/${workspaceId}/email_settings`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  saveEmailSettings: async (
    workspaceId: string,
    payload: SaveEmailSettingsRequest
  ): Promise<SaveEmailSettingsResponse> => {
    try {
      return await apiClient.put<SaveEmailSettingsResponse>(
        `/workspaces/${workspaceId}/email_settings`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
