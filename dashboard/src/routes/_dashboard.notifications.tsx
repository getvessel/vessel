import { createFileRoute } from '@tanstack/react-router';
import { Bell, Plus, Send } from 'lucide-react';
import { EmptyPanel } from '#/features/dashboard/empty-panel';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { useGetNotifications } from '#/hooks/useSettings';

export const Route = createFileRoute('/_dashboard/notifications')({
  component: NotificationsPage,
});

function NotificationsPage() {
  const notifications = useGetNotifications();

  return (
    <OperationalPage
      title="Notifications"
      description="Route deployment failures, backup alerts, update notices, and security events to the right channels."
      scope="Events"
      statusLabel={notifications.isError ? 'Channel error' : 'Ready'}
      statusTone={notifications.isError ? 'danger' : 'healthy'}
      primaryAction={{ label: 'Add channel', icon: Plus }}
      metrics={[
        {
          label: 'Channels',
          value: notifications.isLoading ? undefined : '0',
          detail: 'Configured destinations',
          icon: Bell,
        },
        {
          label: 'Delivery',
          value: 'Quiet',
          detail: 'No active alerts',
          icon: Send,
        },
      ]}
    >
      <EmptyPanel
        icon={Bell}
        title="No notification channels"
        description="Add Slack, webhook, or email channels so important deployment and infrastructure events do not stay hidden."
        actionLabel="Add channel"
      />
    </OperationalPage>
  );
}
