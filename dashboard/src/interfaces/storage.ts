import type { BaseResponse } from './base';

export type StorageType = 'minio' | string;
export type StorageStatus = 'running' | 'stopped' | 'error' | 'starting';

export interface Storage {
  id: string;
  projectId: string;
  environmentId: string;
  name: string;
  type: StorageType;
  apiPort: number;
  consolePort: number;
  accessKey: string;
  secretKey?: string;
  bucketName: string;
  volumePath: string;
  containerId: string;
  status: StorageStatus;
  internalDns: string;
  externalDns: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateStorageRequest {
  projectId: string;
  environmentId: string;
  name: string;
  type: StorageType;
}

export type ListStorageResponse = BaseResponse<Storage[]>;
export type GetStorageResponse = BaseResponse<Storage>;
export type CreateStorageResponse = BaseResponse<Storage>;
