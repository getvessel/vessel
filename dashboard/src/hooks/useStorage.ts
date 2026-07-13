import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { storageService } from '#/services/storage';

export const useList = () => {
  return useQuery({
    queryKey: ['storage', 'list'].filter(Boolean),
    queryFn: () => storageService.list(),
  });
};

export const useGet = (id: string) => {
  return useQuery({
    queryKey: ['storage', 'get', id].filter(Boolean),
    queryFn: () => storageService.get(id),
  });
};

export const useCreate = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { payload: Parameters<typeof storageService.create>[0] }) =>
      storageService.create(payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['storage'] });
    },
  });
};

export const useDelete = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => storageService.delete(payload.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['storage'] });
    },
  });
};

export const useStart = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => storageService.start(payload.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['storage'] });
    },
  });
};

export const useStop = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => storageService.stop(payload.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['storage'] });
    },
  });
};
