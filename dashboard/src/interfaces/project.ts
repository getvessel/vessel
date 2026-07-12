import type { Database } from './database';
import type { AppService } from './deployment';

export interface ProjectConfig {
  id: string;
  workspaceId?: string;
  teamId?: string;
  name: string;
  description?: string;
  createdAt: string;
  updatedAt: string;
}

export interface DomainConfig {
  id: string;
  projectId: string;
  domainName: string;
  redirectTo?: string;
  sslCertStatus: string;
  pathPrefix: string;
  createdAt: string;
  updatedAt: string;
}

export interface ServerlessFunctionCode {
  id: string;
  serviceId: string;
  runtime: string;
  codeContent: string;
  createdAt: string;
  updatedAt: string;
}

export interface EnvironmentConfig {
  id: string;
  projectId: string;
  name: string;
  isDefault: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CanvasSummary {
  id: string;
  workspaceId?: string;
  teamId?: string;
  name: string;
  description?: string;
  createdAt: string;
  updatedAt: string;
  environmentsCount: number;
  appsCount: number;
  databasesCount: number;
  storageCount: number;
  onlineServices: number;
  totalServices: number;
  serviceIcons: string[];
  defaultEnvironment?: EnvironmentConfig;
}

export interface EnvironmentCanvas {
  environment: EnvironmentConfig;
  apps: AppService[];
  databases: Database[];
  storage: Storage[];
}

export interface CreateProjectRequest {
  id: string;
  teamId?: string;
  name: string;
  description?: string;
  repositoryUrl?: string;
  repository_url?: string;
  branch?: string;
  internalPort?: number;
  internal_port?: number;
  domain?: string;
}

export interface ProjectToken {
  id: string;
  projectId: string;
  environmentId: string;
  name: string;
  tokenPrefix: string;
  scopes: string[];
  ipAllowlist: string[];
  expiresAt?: string;
  createdAt: string;
}

export interface ProjectMember {
  id: string;
  projectId: string;
  userId?: string;
  email: string;
  permission: string;
  status: string;
  invitedAt: string;
  acceptedAt?: string;
}

export interface CreateWebhookRequest {
  url: string;
  eventTypes: string[];
  includePrEnvironments: boolean;
}

export interface CreateTokenRequest {
  name: string;
  environmentId: string;
  scopes: string[];
  ipAllowlist?: string[];
  expiresAt?: string;
}

export interface AddMemberRequest {
  email: string;
  permission: string;
}
