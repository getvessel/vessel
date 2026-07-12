import type { BaseEntity } from "./base";

export interface Workspace extends BaseEntity {
  name: string;
  avatarUrl?: string;
  preferredRegion?: string;
  ownerId: string;
}

export interface WorkspaceMember {
  id: string;
  workspaceId: string;
  userId: string;
  userEmail: string;
  role: string;
  joinedAt: string;
}

export interface WorkspaceInvite extends BaseEntity {
  workspaceId: string;
  email: string;
  role: string;
  token: string;
  invitedBy: string;
  expiresAt: string;
}

export interface CreateWorkspaceRequest {
  name: string;
  avatarUrl?: string;
  preferredRegion?: string;
}

export interface InviteMemberRequest {
  email: string;
  role: string;
}

export interface AcceptInviteRequest {
  token: string;
}

export interface GetWorkspaceResponse {
  workspace: Workspace;
  members: WorkspaceMember[];
}

export interface TrustedDomain {
  id: string;
  workspaceId: string;
  domain: string;
  role: string;
  createdAt: string;
}

export interface SSHKey {
  id: string;
  workspaceId: string;
  name: string;
  publicKey: string;
  createdAt: string;
}

export interface AuditLog {
  id: string;
  workspaceId: string;
  userId?: string;
  projectId?: string;
  environmentId?: string;
  action: string;
  resource?: string;
  actor: string;
  ipAddress?: string;
  createdAt: string;
  timestamp?: string;
}

export interface UpdateWorkspaceRequest {
  name?: string;
  avatarUrl?: string;
  preferredRegion?: string;
}

export interface CreateTrustedDomainRequest {
  domain: string;
  role: string;
}

export interface CreateSSHKeyRequest {
  name: string;
  publicKey: string;
}
