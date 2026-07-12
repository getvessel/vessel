import type { AppService, CreateAppServiceRequest, UpdateAppServiceRequest } from '#/interfaces/deployment';
import { apiClient } from './instance';

export const appsService = {
  listByProject: async (projectId: string): Promise<AppService[]> => {
    const { data } = await apiClient.get<AppService[]>(`/projects/${projectId}/apps`);
    return data;
  },

  listByEnvironment: async (environmentId: string): Promise<AppService[]> => {
    const { data } = await apiClient.get<AppService[]>(`/environments/${environmentId}/apps`);
    return data;
  },

  getApp: async (appId: string): Promise<AppService> => {
    const { data } = await apiClient.get<AppService>(`/apps/${appId}`);
    return data;
  },

  createApp: async (environmentId: string, payload: CreateAppServiceRequest): Promise<AppService> => {
    const { data } = await apiClient.post<AppService>(`/environments/${environmentId}/apps`, payload);
    return data;
  },

  updateApp: async (appId: string, payload: UpdateAppServiceRequest): Promise<AppService> => {
    const { data } = await apiClient.put<AppService>(`/apps/${appId}`, payload);
    return data;
  },

  deleteApp: async (appId: string): Promise<void> => {
    await apiClient.delete(`/apps/${appId}`);
  },
};
