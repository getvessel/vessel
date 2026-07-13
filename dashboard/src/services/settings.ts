import type { BaseResponse } from '#/interfaces/base';
import type { ServerSettings, UpdateSettingsRequest } from '#/interfaces/settings';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const settingsService = {
  getSettings: async (): Promise<BaseResponse<ServerSettings>> => {
    try {
      return await apiClient.get<BaseResponse<ServerSettings>>('/settings');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  updateSettings: async (payload: UpdateSettingsRequest): Promise<BaseResponse<ServerSettings>> => {
    try {
      return await apiClient.put<BaseResponse<ServerSettings>>('/settings', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getNotifications: async (): Promise<BaseResponse<unknown>> => {
    try {
      return await apiClient.get<BaseResponse<unknown>>('/settings/notifications');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getOauthProviders: async (): Promise<BaseResponse<unknown>> => {
    try {
      return await apiClient.get<BaseResponse<unknown>>('/settings/oauth/providers');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getGitApps: async (provider: string): Promise<BaseResponse<unknown>> => {
    try {
      return await apiClient.get<BaseResponse<unknown>>(`/settings/git_apps/${provider}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
