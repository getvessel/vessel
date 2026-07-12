import type {
  CreateWorkspaceRequest,
  CreateWorkspaceResponse,
  GetWorkspaceResponse,
  ListWorkspacesResponse,
  UpdateWorkspaceRequest,
  UpdateWorkspaceResponse,
} from '#/interfaces/workspace';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const workspacesService = {
  listWorkspaces: async (): Promise<ListWorkspacesResponse> => {
    try {
      return await apiClient.get<ListWorkspacesResponse>('/workspaces');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getWorkspace: async (id: string): Promise<GetWorkspaceResponse> => {
    try {
      return await apiClient.get<GetWorkspaceResponse>(`/workspaces/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  createWorkspace: async (payload: CreateWorkspaceRequest): Promise<CreateWorkspaceResponse> => {
    try {
      return await apiClient.post<CreateWorkspaceResponse>('/workspaces', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  updateWorkspace: async (
    id: string,
    payload: UpdateWorkspaceRequest
  ): Promise<UpdateWorkspaceResponse> => {
    try {
      return await apiClient.put<UpdateWorkspaceResponse>(`/workspaces/${id}`, payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteWorkspace: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/workspaces/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
