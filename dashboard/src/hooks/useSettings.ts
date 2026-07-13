import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { settingsService } from '#/services/settings';

export const useGetSettings = () => {
  return useQuery({
    queryKey: ['settings', 'getSettings'].filter(Boolean),
    queryFn: () => settingsService.getSettings(),
  });
};

export const useUpdateSettings = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { payload: Parameters<typeof settingsService.updateSettings>[0] }) =>
      settingsService.updateSettings(payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['settings'] });
    },
  });
};

export const useGetNotifications = () => {
  return useQuery({
    queryKey: ['settings', 'getNotifications'].filter(Boolean),
    queryFn: () => settingsService.getNotifications(),
  });
};

export const useGetOauthProviders = () => {
  return useQuery({
    queryKey: ['settings', 'getOauthProviders'].filter(Boolean),
    queryFn: () => settingsService.getOauthProviders(),
  });
};

export const useGetGitApps = (provider: string) => {
  return useQuery({
    queryKey: ['settings', 'getGitApps', provider].filter(Boolean),
    queryFn: () => settingsService.getGitApps(provider),
  });
};
