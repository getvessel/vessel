import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { projectSettingsService } from '#/services/project-settings';

export const useListWebhooks = (projectId: string) => {
  return useQuery({
    queryKey: ['projectSettings', 'listWebhooks', projectId].filter(Boolean),
    queryFn: () => projectSettingsService.listWebhooks(projectId),
  });
};

export const useCreateWebhook = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      projectId: string;
      payload: Parameters<typeof projectSettingsService.createWebhook>[1];
    }) => projectSettingsService.createWebhook(payload.projectId, payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projectSettings'] });
    },
  });
};

export const useDeleteWebhook = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { projectId: string; id: string }) =>
      projectSettingsService.deleteWebhook(payload.projectId, payload.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projectSettings'] });
    },
  });
};

export const useListTokens = (projectId: string) => {
  return useQuery({
    queryKey: ['projectSettings', 'listTokens', projectId].filter(Boolean),
    queryFn: () => projectSettingsService.listTokens(projectId),
  });
};

export const useCreateToken = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      projectId: string;
      payload: Parameters<typeof projectSettingsService.createToken>[1];
    }) => projectSettingsService.createToken(payload.projectId, payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projectSettings'] });
    },
  });
};

export const useDeleteToken = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { projectId: string; id: string }) =>
      projectSettingsService.deleteToken(payload.projectId, payload.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projectSettings'] });
    },
  });
};

export const useListMembers = (projectId: string) => {
  return useQuery({
    queryKey: ['projectSettings', 'listMembers', projectId].filter(Boolean),
    queryFn: () => projectSettingsService.listMembers(projectId),
  });
};

export const useAddMember = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      projectId: string;
      payload: Parameters<typeof projectSettingsService.addMember>[1];
    }) => projectSettingsService.addMember(payload.projectId, payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projectSettings'] });
    },
  });
};

export const useRemoveMember = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { projectId: string; memberId: string }) =>
      projectSettingsService.removeMember(payload.projectId, payload.memberId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projectSettings'] });
    },
  });
};
