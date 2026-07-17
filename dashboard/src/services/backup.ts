import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

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

export interface CreateS3DestinationRequest {
  projectId: string;
  name: string;
  endpoint: string;
  bucket: string;
  region: string;
  accessKeyId: string;
  secretAccessKey: string;
}

export const backupService = {
  listS3Destinations: async (): Promise<{ data: S3Destination[] }> => {
    try {
      return await apiClient.get<{ data: S3Destination[] }>('/s3-destinations');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  createS3Destination: async (
    payload: CreateS3DestinationRequest
  ): Promise<{ data: S3Destination }> => {
    try {
      return await apiClient.post<{ data: S3Destination }>('/s3-destinations', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteS3Destination: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/s3-destinations/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
