export interface TeamNotificationChannel {
  id: string;
  teamId: string;
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
  teamId: string;
  projectId?: string;
  url?: string;
}
