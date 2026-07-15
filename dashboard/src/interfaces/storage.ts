import type { BaseResponse } from './base';

export type StorageVolumeStatus = 'running' | 'stopped' | 'error';

export interface StorageVolume {
  id: string;
  projectId: string;
  environmentId?: string;
  name: string;
  sizeGb: number;
  status: StorageVolumeStatus;
  mountPath: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateStorageRequest {
  projectId: string;
  name: string;
  sizeGb: number;
  mountPath: string;
  environmentId?: string;
}

export type ListStorageResponse = BaseResponse<StorageVolume[]>;
export type GetStorageResponse = BaseResponse<StorageVolume>;
export type CreateStorageResponse = BaseResponse<StorageVolume>;
