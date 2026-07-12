export interface Deployment {
  id: string;
  serviceId: string;
  environmentId: string;
  projectId: string;
  status: string;
  branch?: string;
  commitHash?: string;
  commitMessage?: string;
  trigger?: string;
  buildLogs?: string;
  containerId?: string;
  createdAt: string;
  updatedAt: string;
  finishedAt?: string;
}

export interface ServiceMetric {
  timestamp: string;
  cpuPercent: number;
  memoryMB: number;
  networkRxKB: number;
  networkTxKB: number;
}

export interface TriggerDeploymentRequest {
  branch?: string;
}

export interface AppService {
  id: string;
  projectId: string;
  environmentId: string;
  name: string;
  repositoryUrl: string;
  branch: string;
  rootDirectory: string;
  buildCommand: string;
  startCommand: string;
  dockerfilePath: string;
  buildEngine: string;
  internalPort: number;
  domain: string;
  healthCheckPath: string;
  containerId: string;
  status: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateAppServiceRequest {
  projectId: string;
  name: string;
  repositoryUrl: string;
  branch: string;
  rootDirectory: string;
  buildCommand: string;
  startCommand: string;
  dockerfilePath: string;
  buildEngine: string;
  internalPort: number;
  domain: string;
  healthCheckPath: string;
}

export interface UpdateAppServiceRequest {
  name: string;
  repositoryUrl: string;
  branch: string;
  rootDirectory: string;
  buildCommand: string;
  startCommand: string;
  dockerfilePath: string;
  buildEngine: string;
  internalPort: number;
  domain: string;
  healthCheckPath: string;
  containerId: string;
  status: string;
}

export interface Variable {
  id: string;
  serviceId: string;
  projectId: string;
  environmentId: string;
  key: string;
  value: string;
  isSecret: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreateServiceVarRequest {
  key: string;
  value: string;
  isSecret: boolean;
}

export interface UpdateServiceVarRequest {
  key: string;
  value: string;
  isSecret: boolean;
}

export interface Job {
  id: string;
  projectId: string;
  name: string;
  schedule: string;
  command: string;
  status: string;
  lastRunAt: string;
  lastOutput: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateJobRequest {
  projectId: string;
  name: string;
  schedule: string;
  command: string;
}

export interface UpdateJobRequest {
  name?: string;
  schedule?: string;
  command?: string;
  status?: string;
}

export interface BackupConfig {
  id: string;
  projectId: string;
  databaseId?: string;
  storageId?: string;
  s3DestinationId?: string;
  name: string;
  schedule: string;
  retentionDays: number;
  status: string;
  createdAt: string;
  updatedAt: string;
}

export interface BackupRecord {
  id: string;
  backupConfigId: string;
  projectId: string;
  databaseId?: string;
  status: string;
  filePath: string;
  fileSizeBytes: number;
  s3Url?: string;
  logs: string;
  startedAt: string;
  completedAt: string;
}

export interface S3Destination {
  id: string;
  projectId: string;
  name: string;
  endpoint: string;
  bucket: string;
  region: string;
  accessKeyId: string;
  secretAccessKey: string;
  createdAt: string;
}

export interface PRPreview {
  id: string;
  serviceId: string;
  projectId: string;
  prNumber: number;
  branch: string;
  commitHash: string;
  status: string;
  previewDomain: string;
  containerId: string;
  createdAt: string;
  updatedAt: string;
}
