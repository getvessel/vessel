import type { BaseResponse } from '#/interfaces/base';
import type {
  CreateDatabaseRequest,
  CreateDatabaseResponse,
  DatabaseQueryResponseType,
  DeleteDatabaseResponse,
  GetDatabaseResponse,
  GetDatabasesResponse,
  ImportDatabaseRequest,
  ListTablesResponse,
  TableRowPayload,
} from '#/interfaces/database';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const databasesService = {
  getDatabases: async (projectId: string): Promise<GetDatabasesResponse> => {
    try {
      return await apiClient.get<GetDatabasesResponse>(`/projects/${projectId}/databases`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getDatabase: async (id: string): Promise<GetDatabaseResponse> => {
    try {
      return await apiClient.get<GetDatabaseResponse>(`/databases/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  createDatabase: async (payload: CreateDatabaseRequest): Promise<CreateDatabaseResponse> => {
    try {
      return await apiClient.post<CreateDatabaseResponse>('/databases', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteDatabase: async (id: string): Promise<DeleteDatabaseResponse> => {
    try {
      return await apiClient.delete<DeleteDatabaseResponse>(`/databases/${id}`);
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

  executeQuery: async (id: string, query: string): Promise<DatabaseQueryResponseType> => {
    try {
      return await apiClient.post<DatabaseQueryResponseType>(`/databases/${id}/query`, { query });
    } catch (error) {
      throw handleApiError(error);
    }
  },

  listTables: async (id: string): Promise<ListTablesResponse> => {
    try {
      return await apiClient.get<ListTablesResponse>(`/databases/${id}/tables`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getTableData: async (
    id: string,
    table: string,
    limit: number = 100,
    offset: number = 0
  ): Promise<DatabaseQueryResponseType> => {
    try {
      return await apiClient.get<DatabaseQueryResponseType>(
        `/databases/${id}/data/${table}?limit=${limit}&offset=${offset}`
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  insertTableRow: async (
    id: string,
    table: string,
    payload: TableRowPayload
  ): Promise<DatabaseQueryResponseType> => {
    try {
      return await apiClient.post<DatabaseQueryResponseType>(
        `/databases/${id}/data/${table}`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  updateTableRow: async (
    id: string,
    table: string,
    payload: TableRowPayload
  ): Promise<DatabaseQueryResponseType> => {
    try {
      return await apiClient.put<DatabaseQueryResponseType>(
        `/databases/${id}/data/${table}`,
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteTableRow: async (id: string, table: string, payload?: TableRowPayload): Promise<void> => {
    try {
      await apiClient.delete(`/databases/${id}/data/${table}`, {
        body: payload ? JSON.stringify(payload) : undefined,
      });
    } catch (error) {
      throw handleApiError(error);
    }
  },

  importDatabase: async (
    id: string,
    payload: ImportDatabaseRequest
  ): Promise<BaseResponse<void>> => {
    try {
      return await apiClient.post<BaseResponse<void>>(`/databases/${id}/import`, payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
