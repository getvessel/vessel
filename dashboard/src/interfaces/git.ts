export interface GitProviderConfig {
  id: string;
  userId: string;
  provider: string;
  accessToken?: string;
  accountName: string;
  createdAt: string;
  updatedAt: string;
}

export interface GitRepository {
  id: number;
  name: string;
  fullName: string;
  private: boolean;
  cloneUrl: string;
  htmlUrl: string;
  defaultBranch: string;
}

export interface GitConnectRequest {
  provider: string;
  accessToken: string;
  accountName: string;
}

export interface GithubApp {
  id: string;
  workspaceId: string;
  name: string;
  appId: string;
  installationId: string;
  clientId: string;
  clientSecret: string;
  webhookSecret: string;
  privateKey: string;
  isPublic: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface GitlabApp {
  id: string;
  workspaceId: string;
  name: string;
  appId: string;
  appSecret: string;
  webhookSecret: string;
  apiUrl: string;
  isPublic: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface BitbucketApp {
  id: string;
  workspaceId: string;
  name: string;
  workspace: string;
  clientId: string;
  clientSecret: string;
  webhookSecret: string;
  isPublic: boolean;
  createdAt: string;
  updatedAt: string;
}
