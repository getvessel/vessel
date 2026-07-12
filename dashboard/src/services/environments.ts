import type { EnvironmentConfig } from '#/interfaces/project';
import { apiClient } from './instance';

export const environmentsService = {
  listByProject: async (projectId: string): Promise<EnvironmentConfig[]> => {
    const { data } = await apiClient.get<EnvironmentConfig[]>(`/projects/${projectId}/environments`);
    return data;
  },

  createEnvironment: async (projectId: string, name: string): Promise<EnvironmentConfig> => {
    const { data } = await apiClient.post<EnvironmentConfig>(`/projects/${projectId}/environments`, { name });
    return data;
  },

  deleteEnvironment: async (environmentId: string): Promise<void> => {
    await apiClient.delete(`/environments/${environmentId}`);
  },
};
