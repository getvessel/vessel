import { createFileRoute } from '@tanstack/react-router';
import { ArchiveRestore, DatabaseBackup, Download, ShieldCheck } from 'lucide-react';
import { EmptyPanel } from '#/features/dashboard/empty-panel';
import { OperationalPage } from '#/features/dashboard/operational-page';

export const Route = createFileRoute('/_dashboard/databases/$dbId/backups')({
  component: DatabaseBackupsPage,
});

function DatabaseBackupsPage() {
  const { dbId } = Route.useParams();

  return (
    <OperationalPage
      title="Database backups"
      description="Schedule snapshots, verify retention, and restore a known-good copy when a database needs recovery."
      scope={dbId}
      statusLabel="Policy ready"
      statusTone="healthy"
      primaryAction={{ label: 'Create backup', icon: DatabaseBackup }}
      secondaryActions={[{ label: 'Restore', icon: ArchiveRestore, variant: 'outline' }]}
      metrics={[
        {
          label: 'Retention',
          value: '7 days',
          detail: 'Default policy',
          icon: ShieldCheck,
        },
        {
          label: 'Latest backup',
          value: 'None',
          detail: 'Create the first snapshot',
          icon: Download,
        },
      ]}
    >
      <EmptyPanel
        icon={DatabaseBackup}
        title="No backups recorded"
        description="Create an on-demand backup or enable scheduled backups to protect this database before risky changes."
        actionLabel="Create backup"
      />
    </OperationalPage>
  );
}
