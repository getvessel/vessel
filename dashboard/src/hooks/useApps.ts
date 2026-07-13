import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { appsService } from '#/services/apps';

export const useListByProject = (projectId: string) => {
  return useQuery({
    queryKey: ['apps', 'listByProject', projectId].filter(Boolean),
    queryFn: () => appsService.listByProject(projectId),
  });
};

export const useListByEnvironment = (environmentId: string) => {
  return useQuery({
    queryKey: ['apps', 'listByEnvironment', environmentId].filter(Boolean),
    queryFn: () => appsService.listByEnvironment(environmentId),
  });
};

export const useGetApp = (appId: string) => {
  return useQuery({
    queryKey: ['apps', 'getApp', appId].filter(Boolean),
    queryFn: () => appsService.getApp(appId),
  });
};

export const useCreateApp = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      environmentId: string;
      payload: Parameters<typeof appsService.createApp>[1];
    }) => appsService.createApp(payload.environmentId, payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['apps'] });
    },
  });
};

export const useUpdateApp = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      appId: string;
      payload: Parameters<typeof appsService.updateApp>[1];
    }) => appsService.updateApp(payload.appId, payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['apps'] });
    },
  });
};

export const useDeleteApp = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { appId: string }) => appsService.deleteApp(payload.appId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['apps'] });
    },
  });
};
