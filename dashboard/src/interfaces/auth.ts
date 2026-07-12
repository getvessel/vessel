import type { User } from './users';

export interface AuthResponse {
  token: string;
  user: User;
}

export interface ApiErrorResponse {
  error: string;
}
