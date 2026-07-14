export interface UserVercelAccount {
  id: string;
  userId: string;
  accessToken?: string;
  accountName: string;
  createdAt: string;
  updatedAt: string;
}

export interface VercelProject {
  id: string;
  name: string;
  framework: string;
  nodeVersion: string;
  accountId: string;
}

export interface VercelEnvVar {
  type: string;
  key: string;
  value: string;
  target: string[];
}
