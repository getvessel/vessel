import { createFileRoute } from '@tanstack/react-router';
import { HardDrive, Save, Trash2, Wrench } from 'lucide-react';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { SettingsPanel } from '#/features/dashboard/settings-panel';

export const Route = createFileRoute('/_dashboard/settings/maintenance')({
  component: MaintenanceSettingsPage,
});

function MaintenanceSettingsPage() {
  return (
    <OperationalPage
      title="Maintenance"
      description="Control cleanup windows, prune old artifacts, and keep node-level maintenance actions explicit."
      scope="Operations"
      statusLabel="Scheduled"
      statusTone="healthy"
      primaryAction={{ label: 'Save window', icon: Save }}
      secondaryActions={[{ label: 'Run cleanup', icon: Trash2, variant: 'outline' }]}
      metrics={[
        {
          label: 'Window',
          value: '02:00',
          detail: 'Local maintenance time',
          icon: Wrench,
        },
        {
          label: 'Artifacts',
          value: 'Pruned',
          detail: 'Old build cache',
          icon: HardDrive,
        },
      ]}
    >
      <SettingsPanel
        title="Cleanup policy"
        icon={Wrench}
        fields={[
          { label: 'Maintenance window', value: '02:00-03:00' },
          { label: 'Keep deployment logs', value: '30 days' },
          { label: 'Keep build cache', value: '14 days' },
        ]}
        toggles={[
          {
            label: 'Auto prune unused images',
            description:
              'Remove images and containers that are no longer referenced by active services.',
            checked: true,
          },
        ]}
      />
    </OperationalPage>
  );
}
