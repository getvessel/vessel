import { useQuery } from '@tanstack/react-query';
import { vercelService } from '#/services/vercel';

export const useListProjects = () => {
  return useQuery({
    queryKey: ['vercel', 'listProjects'].filter(Boolean),
    queryFn: () => vercelService.listProjects(),
  });
};

export const useGetProjectEnv = (id: string) => {
  return useQuery({
    queryKey: ['vercel', 'getProjectEnv', id].filter(Boolean),
    queryFn: () => vercelService.getProjectEnv(id),
  });
};
