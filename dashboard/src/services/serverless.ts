import type { BaseResponse } from '#/interfaces/base';
import type { ServerlessFunctionCode } from '#/interfaces/project';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const serverlessService = {
  getCode: async (serviceId: string): Promise<BaseResponse<ServerlessFunctionCode>> => {
    try {
      return await apiClient.get<BaseResponse<ServerlessFunctionCode>>(
        `/services/${serviceId}/serverless/code`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  saveCode: async (
    serviceId: string,
    payload: { codeContent: string; runtime?: string }
  ): Promise<BaseResponse<ServerlessFunctionCode>> => {
    try {
      return await apiClient.post<BaseResponse<ServerlessFunctionCode>>(
        `/services/${serviceId}/serverless/code`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
