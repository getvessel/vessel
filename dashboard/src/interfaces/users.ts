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
