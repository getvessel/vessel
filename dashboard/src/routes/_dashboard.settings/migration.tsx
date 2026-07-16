import { createFileRoute } from '@tanstack/react-router';
import { ArrowRightLeft, Download, UploadCloud } from 'lucide-react';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { SettingsPanel } from '#/features/dashboard/settings-panel';

export const Route = createFileRoute('/_dashboard/settings/migration')({
  component: MigrationSettingsPage,
});

function MigrationSettingsPage() {
  return (
    <OperationalPage
      title="Migration"
      description="Export instance metadata, prepare imports, and validate compatibility before moving workloads."
      scope="Portability"
      statusLabel="Ready"
      statusTone="info"
      primaryAction={{ label: 'Export bundle', icon: Download }}
      secondaryActions={[{ label: 'Import bundle', icon: UploadCloud, variant: 'outline' }]}
      metrics={[
        {
          label: 'Mode',
          value: 'Metadata',
          detail: 'Databases handled by backup',
          icon: ArrowRightLeft,
        },
      ]}
    >
      <SettingsPanel
        title="Migration bundle"
        icon={ArrowRightLeft}
        fields={[
          { label: 'Export scope', value: 'projects, services, settings' },
          { label: 'Encryption recipient', value: 'admin@example.com' },
        ]}
        toggles={[
          {
            label: 'Validate after import',
            description: 'Run schema and dependency checks before enabling imported workloads.',
            checked: true,
          },
        ]}
      />
    </OperationalPage>
  );
}
