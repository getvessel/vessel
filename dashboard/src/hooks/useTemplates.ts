import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { templatesService } from '#/services/templates';
import type { OneClickDeployRequest } from '#/interfaces/templates';

export const useListOneClickApps = () => {
  return useQuery({
    queryKey: ['oneClickApps'],
    queryFn: () => templatesService.listOneClickApps(),
  });
};

export const useDeployOneClickApp = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: OneClickDeployRequest) => templatesService.deployOneClickApp(payload),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['projects', variables.projectId, 'apps'] });
    },
  });
};

export const useDeployCompose = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ projectId, file }: { projectId: string; file: File }) =>
      templatesService.deployCompose(projectId, file),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['projects', variables.projectId, 'apps'] });
    },
  });
};

export const useDeployArchive = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ projectId, file }: { projectId: string; file: File }) =>
      templatesService.deployArchive(projectId, file),
    onSuccess: (_, variables) => {
      if (variables.projectId) {
        queryClient.invalidateQueries({ queryKey: ['projects', variables.projectId, 'apps'] });
      }
    },
  });
};
