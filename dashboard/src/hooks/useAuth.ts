import { useMutation, useQueryClient } from '@tanstack/react-query';
import { authService } from '#/services/auth';
import { authActions } from '#/stores/authStore';

export const useLogin = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { credentials: Parameters<typeof authService.login>[0] }) =>
      authService.login(payload.credentials),
    onSuccess: (data) => {
      if (data?.token && data?.user) {
        authActions.setAuth(data.token, data.user);
      }
      queryClient.invalidateQueries({ queryKey: ['auth'] });
    },
  });
};

export const useRegister = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { details: Parameters<typeof authService.register>[0] }) =>
      authService.register(payload.details),
    onSuccess: (data) => {
      if (data?.token && data?.user) {
        authActions.setAuth(data.token, data.user);
      }
      queryClient.invalidateQueries({ queryKey: ['auth'] });
    },
  });
};

export const useLogout = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => authService.logout(),
    onSuccess: () => {
      authActions.logout();
      queryClient.clear();
    },
  });
};

export const useSetup2FA = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => authService.setup2FA(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['auth'] });
    },
  });
};

export const useVerify2FA = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { payload: Parameters<typeof authService.verify2FA>[0] }) =>
      authService.verify2FA(payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['auth'] });
    },
  });
};

export const useDisable2FA = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { payload: Parameters<typeof authService.disable2FA>[0] }) =>
      authService.disable2FA(payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['auth'] });
    },
  });
};
