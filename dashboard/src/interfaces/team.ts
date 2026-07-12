export interface Team {
  id: string;
  name: string;
  avatarUrl?: string;
  preferredRegion?: string;
  ownerId: string;
  createdAt: string;
  updatedAt: string;
}

export interface TeamMember {
  id: string;
  teamId: string;
  userId: string;
  userEmail: string;
  role: string;
  joinedAt: string;
}

export interface TeamInvite {
  id: string;
  teamId: string;
  email: string;
  role: string;
  token: string;
  invitedBy: string;
  expiresAt: string;
  createdAt: string;
}

export interface CreateTeamRequest {
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

export interface GetTeamResponse {
  team: Team;
  members: TeamMember[];
}

export interface Workspace {
  id: string;
  name: string;
  avatarUrl?: string;
  preferredRegion?: string;
  ownerId: string;
  createdAt: string;
  updatedAt: string;
}

export interface TrustedDomain {
  id: string;
  teamId: string;
  domain: string;
  role: string;
  createdAt: string;
}

export interface SSHKey {
  id: string;
  teamId: string;
  name: string;
  publicKey: string;
  createdAt: string;
}

export interface AuditLog {
  id: string;
  teamId: string;
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

export interface CreateWorkspaceRequest {
  name: string;
  avatarUrl?: string;
  preferredRegion?: string;
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
