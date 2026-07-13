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
  email?: string;
  avatarUrl?: string;
}

export interface ChangePasswordRequest {
  oldPassword?: string;
  newPassword?: string;
}

export type GetProfileResponse = BaseResponse<UserProfile>;
export type UpdateProfileResponse = BaseResponse<UserProfile>;
