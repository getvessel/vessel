import type { BaseResponse, PaginatedData } from './base';
import type { Database } from './database';
import type { AppService } from './deployment';

export interface ProjectConfig {
  id: string;
  workspaceId?: string;
  name: string;
  description?: string;
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

export interface CanvasSummary {
  id: string;
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
  createdAt: string;
  updatedAt: string;
}

export interface EnvironmentCanvas {
  environment: EnvironmentConfig;
  apps: AppService[];
  databases: Database[];
  storage: Storage[];
}

export interface CreateProjectRequest {
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

export type ListProjectsResponse = BaseResponse<PaginatedData<ProjectConfig>>;
export type GetProjectResponse = BaseResponse<ProjectConfig>;
export type CreateProjectResponse = BaseResponse<ProjectConfig>;
export type ListEnvironmentsResponse = BaseResponse<EnvironmentConfig[]>;
export type CreateEnvironmentResponse = BaseResponse<EnvironmentConfig>;
