import { Clock3, Plus } from 'lucide-react';
import { EmptyPanel } from '#/features/dashboard/empty-panel';
import { OperationalPage } from '#/features/dashboard/operational-page';

export const JobList = () => (
  <OperationalPage
    title="Jobs"
    description="Create and monitor scheduled commands for maintenance and background work."
    scope="Scheduler"
    statusLabel="Ready"
    statusTone="healthy"
    primaryAction={{ label: 'New job', icon: Plus }}
    metrics={[
      {
        label: 'Jobs',
        value: '0',
        detail: 'Scheduled tasks',
        icon: Clock3,
      },
    ]}
  >
    <EmptyPanel
      icon={Clock3}
      title="No jobs configured"
      description="Add a recurring command to run cleanup, synchronization, or background maintenance."
      actionLabel="New job"
    />
  </OperationalPage>
);
