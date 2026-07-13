import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { workspaceSettingsService } from '#/services/workspace-settings';

export const useGetAISettings = (workspaceId: string) => {
  return useQuery({
    queryKey: ['workspaceSettings', 'getAISettings', workspaceId].filter(Boolean),
    queryFn: () => workspaceSettingsService.getAISettings(workspaceId),
  });
};

export const useSaveAISettings = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      workspaceId: string;
      payload: Parameters<typeof workspaceSettingsService.saveAISettings>[1];
    }) => workspaceSettingsService.saveAISettings(payload.workspaceId, payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['workspaceSettings'] });
    },
  });
};

export const useGetEmailSettings = (workspaceId: string) => {
  return useQuery({
    queryKey: ['workspaceSettings', 'getEmailSettings', workspaceId].filter(Boolean),
    queryFn: () => workspaceSettingsService.getEmailSettings(workspaceId),
  });
};

export const useSaveEmailSettings = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      workspaceId: string;
      payload: Parameters<typeof workspaceSettingsService.saveEmailSettings>[1];
    }) => workspaceSettingsService.saveEmailSettings(payload.workspaceId, payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['workspaceSettings'] });
    },
  });
};
