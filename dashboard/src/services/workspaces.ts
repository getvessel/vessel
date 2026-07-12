import type { CreateWorkspaceRequest, GetWorkspaceResponse, UpdateWorkspaceRequest, Workspace } from '#/interfaces/workspace';
import { apiClient } from './instance';

export const workspacesService = {
  listWorkspaces: async (): Promise<Workspace[]> => {
    const { data } = await apiClient.get<Workspace[]>('/workspaces');
    return data;
  },

  getWorkspace: async (id: string): Promise<GetWorkspaceResponse> => {
    const { data } = await apiClient.get<GetWorkspaceResponse>(`/workspaces/${id}`);
    return data;
  },

  createWorkspace: async (payload: CreateWorkspaceRequest): Promise<Workspace> => {
    const { data } = await apiClient.post<Workspace>('/workspaces', payload);
    return data;
  },

  updateWorkspace: async (id: string, payload: UpdateWorkspaceRequest): Promise<Workspace> => {
    const { data } = await apiClient.put<Workspace>(`/workspaces/${id}`, payload);
    return data;
  },

  deleteWorkspace: async (id: string): Promise<void> => {
    await apiClient.delete(`/workspaces/${id}`);
  },
};
