import type { BaseResponse } from '#/interfaces/base';
import type { DomainConfig } from '#/interfaces/project';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const domainsService = {
  listByProject: async (projectId: string): Promise<BaseResponse<DomainConfig[]>> => {
    try {
      return await apiClient.get<BaseResponse<DomainConfig[]>>(`/projects/${projectId}/domains`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  create: async (
    projectId: string,
    payload: { domainName: string; redirectTo?: string; pathPrefix?: string }
  ): Promise<BaseResponse<DomainConfig>> => {
    try {
      return await apiClient.post<BaseResponse<DomainConfig>>(
        `/projects/${projectId}/domains`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  delete: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/domains/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
