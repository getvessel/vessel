import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { backupService } from '#/services/backup';

export function useS3Destinations() {
  return useQuery({
    queryKey: ['s3-destinations'],
    queryFn: () => backupService.listS3Destinations(),
  });
}

export function useCreateS3Destination() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: backupService.createS3Destination,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['s3-destinations'] });
    },
  });
}

export function useDeleteS3Destination() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: backupService.deleteS3Destination,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['s3-destinations'] });
    },
  });
}
