import type { BaseResponse } from '#/interfaces/base';
import type {
  CreateDatabaseRequest,
  Database,
  DatabaseQueryRequest,
  DatabaseQueryResponse,
} from '#/interfaces/database';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const databasesService = {
  listDatabases: async (): Promise<BaseResponse<Database[]>> => {
    try {
      return await apiClient.get<BaseResponse<Database[]>>('/databases');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getDatabase: async (id: string): Promise<BaseResponse<Database>> => {
    try {
      return await apiClient.get<BaseResponse<Database>>(`/databases/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  createDatabase: async (payload: CreateDatabaseRequest): Promise<BaseResponse<Database>> => {
    try {
      return await apiClient.post<BaseResponse<Database>>('/databases', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteDatabase: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/databases/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  startDatabase: async (id: string): Promise<void> => {
    try {
      await apiClient.post(`/databases/${id}/start`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  stopDatabase: async (id: string): Promise<void> => {
    try {
      await apiClient.post(`/databases/${id}/stop`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  queryDatabase: async (
    id: string,
    payload: DatabaseQueryRequest
  ): Promise<BaseResponse<DatabaseQueryResponse>> => {
    try {
      return await apiClient.post<BaseResponse<DatabaseQueryResponse>>(
        `/databases/${id}/query`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
