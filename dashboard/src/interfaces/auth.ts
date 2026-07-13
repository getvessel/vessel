import type { User } from './users';

export interface AuthResponse {
  token: string;
  user: User;
}

export interface ApiErrorResponse {
  error: string;
}

export interface AuthCredentials {
  email: string;
  password: string;
}

export interface RegisterCredentials extends AuthCredentials {
  name: string;
}

export interface Setup2FAResponse {
  secret: string;
  qrCodeUrl: string;
}

export interface Verify2FARequest {
  token: string;
}
