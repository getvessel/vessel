import { createFileRoute } from '@tanstack/react-router';
import { NotificationsSettings } from '#/features/instance/notifications-settings';

export const Route = createFileRoute('/_dashboard/settings/notifications')({
  component: NotificationsSettings,
});
