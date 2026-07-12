import type { BaseEntity } from './base';
import type { Database } from './database';
import type { AppService } from './deployment';

export interface ProjectConfig extends BaseEntity {
  workspaceId?: string;
  name: string;
  description?: string;
}

export interface DomainConfig extends BaseEntity {
  projectId: string;
  domainName: string;
  redirectTo?: string;
  sslCertStatus: string;
  pathPrefix: string;
}

export interface ServerlessFunctionCode extends BaseEntity {
  serviceId: string;
  runtime: string;
  codeContent: string;
}

export interface EnvironmentConfig extends BaseEntity {
  projectId: string;
  name: string;
  isDefault: boolean;
}

export interface CanvasSummary extends BaseEntity {
  workspaceId?: string;
  name: string;
  description?: string;
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
  id?: string;
  workspaceId?: string;
  name: string;
  description?: string;
  repositoryUrl?: string;
  branch?: string;
  internalPort?: number;
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
