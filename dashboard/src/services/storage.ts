import type {
  CreateStorageRequest,
  CreateStorageResponse,
  GetStorageResponse,
  ListStorageResponse,
} from '#/interfaces/storage';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const storageService = {
  list: async (): Promise<ListStorageResponse> => {
    try {
      return await apiClient.get<ListStorageResponse>('/storage');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  get: async (id: string): Promise<GetStorageResponse> => {
    try {
      return await apiClient.get<GetStorageResponse>(`/storage/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  create: async (payload: CreateStorageRequest): Promise<CreateStorageResponse> => {
    try {
      return await apiClient.post<CreateStorageResponse>('/storage', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  delete: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/storage/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  start: async (id: string): Promise<void> => {
    try {
      await apiClient.post(`/storage/${id}/start`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  stop: async (id: string): Promise<void> => {
    try {
      await apiClient.post(`/storage/${id}/stop`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
