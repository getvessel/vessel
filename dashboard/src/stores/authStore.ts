import { create } from 'zustand';
import type { User } from '#/interfaces';

export interface AuthState {
  token: string | null;
  user: User | null;
  isAuthenticated: boolean;
  setAuth: (token: string, user: User) => void;
  logout: () => void;
}

const isBrowser = typeof window !== 'undefined';

const getInitialAuth = () => {
  if (!isBrowser) {
    return { token: null, user: null, isAuthenticated: false };
  }

  const storedToken = localStorage.getItem('vessl_auth_token');
  const storedUser = localStorage.getItem('vessl_auth_user');
  let parsedUser = null;
  if (storedUser) {
    try {
      parsedUser = JSON.parse(storedUser);
    } catch (e) {
      console.error('Failed to parse stored user:', e);
      localStorage.removeItem('vessl_auth_user');
    }
  }

  return {
    token: storedToken,
    user: parsedUser,
    isAuthenticated: !!storedToken,
  };
};

export const useAuthStore = create<AuthState>((set) => ({
  ...getInitialAuth(),
  setAuth: (token: string, user: User) => {
    set({ token, user, isAuthenticated: true });
  },
  logout: () => {
    set({ token: null, user: null, isAuthenticated: false });
  },
}));

if (isBrowser) {
  useAuthStore.subscribe((state) => {
    if (state.token) {
      localStorage.setItem('vessl_auth_token', state.token);
    } else {
      localStorage.removeItem('vessl_auth_token');
    }

    if (state.user) {
      localStorage.setItem('vessl_auth_user', JSON.stringify(state.user));
    } else {
      localStorage.removeItem('vessl_auth_user');
    }
  });
}
