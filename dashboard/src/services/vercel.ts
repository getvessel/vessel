import type { BaseResponse } from '#/interfaces/base';
import type { VercelEnvVar, VercelProject } from '#/interfaces/vercel';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const vercelService = {
  listProjects: async (): Promise<BaseResponse<VercelProject[]>> => {
    try {
      return await apiClient.get<BaseResponse<VercelProject[]>>('/vercel/projects');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getProjectEnv: async (id: string): Promise<BaseResponse<VercelEnvVar[]>> => {
    try {
      return await apiClient.get<BaseResponse<VercelEnvVar[]>>(`/vercel/projects/${id}/env`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
