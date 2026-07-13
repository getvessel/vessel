import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { profileService } from '#/services/profile';

export const useGetProfile = () => {
  return useQuery({
    queryKey: ['profile', 'getProfile'].filter(Boolean),
    queryFn: () => profileService.getProfile(),
  });
};

export const useUpdateProfile = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { payload: Parameters<typeof profileService.updateProfile>[0] }) =>
      profileService.updateProfile(payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profile'] });
    },
  });
};

export const useListTokens = () => {
  return useQuery({
    queryKey: ['profile', 'listTokens'].filter(Boolean),
    queryFn: () => profileService.listTokens(),
  });
};

export const useCreateToken = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { payload: Parameters<typeof profileService.createToken>[0] }) =>
      profileService.createToken(payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profile'] });
    },
  });
};

export const useDeleteToken = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => profileService.deleteToken(payload.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profile'] });
    },
  });
};
