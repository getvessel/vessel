export interface User {
  id: string;
  email: string;
  name: string;
  role: string;
  totpEnabled: boolean;
  oauthProvider?: string;
  createdAt: string;
  updatedAt: string;
}

export interface PersonalAccessToken {
  id: string;
  userId: string;
  name: string;
  prefix: string;
  expiresAt?: string;
  createdAt: string;
}

export interface CreatePATResponse {
  pat: PersonalAccessToken;
  token: string;
}
