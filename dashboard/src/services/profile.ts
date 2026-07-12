import type { CreatePATResponse, PersonalAccessToken, User } from '#/interfaces/users';
import { apiClient } from './instance';

export const profileService = {
  getProfile: async (): Promise<User> => {
    const { data } = await apiClient.get<User>('/profile');
    return data;
  },

  updateProfile: async (payload: { email?: string; role?: string }): Promise<User> => {
    const { data } = await apiClient.put<User>('/profile', payload);
    return data;
  },

  listPATs: async (): Promise<PersonalAccessToken[]> => {
    const { data } = await apiClient.get<PersonalAccessToken[]>('/profile/tokens');
    return data;
  },

  createPAT: async (payload: { name: string }): Promise<CreatePATResponse> => {
    const { data } = await apiClient.post<CreatePATResponse>('/profile/tokens', payload);
    return data;
  },

  deletePAT: async (id: string): Promise<void> => {
    await apiClient.delete(`/profile/tokens/${id}`);
  },
};
