import type { BaseResponse } from './base';

export interface AISettings {
  workspaceId: string;
  openaiApiKey?: string;
  anthropicApiKey?: string;
  defaultProvider?: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface EmailSettings {
  configured: boolean;
  provider?: string;
  fromAddress?: string;
  fromName?: string;
}

export interface SaveAISettingsRequest {
  openaiApiKey?: string;
  anthropicApiKey?: string;
  defaultProvider?: string;
}

export interface SaveEmailSettingsRequest {
  provider: string;
  fromAddress: string;
  fromName: string;
  apiKey?: string;
  smtpHost?: string;
  smtpPort?: number;
  smtpUser?: string;
  smtpPassword?: string;
}

// Response Types
export type GetAISettingsResponse = BaseResponse<AISettings>;
export type SaveAISettingsResponse = BaseResponse<AISettings>;
export type GetEmailSettingsResponse = BaseResponse<EmailSettings>;
export type SaveEmailSettingsResponse = BaseResponse<EmailSettings>;
