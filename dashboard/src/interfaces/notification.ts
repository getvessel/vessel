export interface TeamNotificationChannel {
  id: string;
  workspaceId: string;
  provider: string;
  config: Record<string, unknown>;
  events: Record<string, unknown>;
  isEnabled: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface NotificationEvent {
  title: string;
  message: string;
  level: string;
  eventType: string;
  workspaceId: string;
  projectId?: string;
  url?: string;
}
