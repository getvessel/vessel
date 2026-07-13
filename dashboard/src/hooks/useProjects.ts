import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { projectsService } from '#/services/projects';

export const useListProjects = (workspaceId: string) => {
  return useQuery({
    queryKey: ['projects', 'listProjects', workspaceId].filter(Boolean),
    queryFn: () => projectsService.listProjects(workspaceId),
  });
};

export const useGetProject = (id: string) => {
  return useQuery({
    queryKey: ['projects', 'getProject', id].filter(Boolean),
    queryFn: () => projectsService.getProject(id),
  });
};

export const useCreateProject = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { payload: Parameters<typeof projectsService.createProject>[0] }) =>
      projectsService.createProject(payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects'] });
    },
  });
};

export const useDeleteProject = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => projectsService.deleteProject(payload.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects'] });
    },
  });
};

export const useDeployProject = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => projectsService.deployProject(payload.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects'] });
    },
  });
};

export const useGetVars = (id: string) => {
  return useQuery({
    queryKey: ['projects', 'getVars', id].filter(Boolean),
    queryFn: () => projectsService.getVars(id),
  });
};

export const useSetVars = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string; payload: { variables: Record<string, string> } }) =>
      projectsService.setVars(payload.id, payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects'] });
    },
  });
};
