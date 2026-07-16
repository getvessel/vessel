import { Plus, Server } from 'lucide-react';
import { EmptyPanel } from '#/features/dashboard/empty-panel';
import { OperationalPage } from '#/features/dashboard/operational-page';

export const StorageBucketList = () => (
  <OperationalPage
    title="Storage buckets"
    description="Manage object storage buckets, access policy, and endpoint health."
    scope="Storage"
    statusLabel="Ready"
    statusTone="info"
    primaryAction={{ label: 'New bucket', icon: Plus }}
    metrics={[
      {
        label: 'Buckets',
        value: '0',
        detail: 'Configured stores',
        icon: Server,
      },
    ]}
  >
    <EmptyPanel
      icon={Server}
      title="No buckets configured"
      description="Create a bucket to store uploads, static assets, backups, or service-generated objects."
      actionLabel="New bucket"
    />
  </OperationalPage>
);
