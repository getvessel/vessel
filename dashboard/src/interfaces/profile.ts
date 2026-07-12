import type { BaseResponse } from './base';

export interface UserProfile {
  id: string;
  email: string;
  name: string;
  role: string;
  avatarUrl?: string;
  totpEnabled: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface UpdateProfileRequest {
  name?: string;
  avatarUrl?: string;
}

// Response Types
export type GetProfileResponse = BaseResponse<UserProfile>;
export type UpdateProfileResponse = BaseResponse<UserProfile>;
