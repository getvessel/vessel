import { createFileRoute } from '@tanstack/react-router';
import { ArchiveRestore, DatabaseBackup, Save } from 'lucide-react';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { SettingsPanel } from '#/features/dashboard/settings-panel';

export const Route = createFileRoute('/_dashboard/settings/backups')({
  component: BackupSettingsPage,
});

function BackupSettingsPage() {
  return (
    <OperationalPage
      title="Backups"
      description="Set instance-wide backup destinations, retention, schedules, and restore safeguards."
      scope="Resilience"
      statusLabel="Policy ready"
      statusTone="healthy"
      primaryAction={{ label: 'Save policy', icon: Save }}
      secondaryActions={[{ label: 'Restore backup', icon: ArchiveRestore, variant: 'outline' }]}
      metrics={[
        {
          label: 'Retention',
          value: '7 days',
          detail: 'Default snapshots',
          icon: DatabaseBackup,
        },
      ]}
    >
      <SettingsPanel
        title="Snapshot policy"
        icon={DatabaseBackup}
        fields={[
          { label: 'Schedule', value: '0 2 * * *' },
          { label: 'Retention window', value: '7 days' },
          { label: 'Destination', value: 'local:/var/lib/vessl/backups' },
        ]}
        toggles={[
          {
            label: 'Verify backup after write',
            description: 'Run integrity checks after each snapshot finishes.',
            checked: true,
          },
        ]}
      />
    </OperationalPage>
  );
}
