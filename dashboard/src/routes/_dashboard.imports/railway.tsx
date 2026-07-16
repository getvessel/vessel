import { createFileRoute } from '@tanstack/react-router';
import { Import, KeyRound, TrainFront, UploadCloud } from 'lucide-react';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { SettingsPanel } from '#/features/dashboard/settings-panel';

export const Route = createFileRoute('/_dashboard/imports/railway')({
  component: RailwayImportPage,
});

function RailwayImportPage() {
  return (
    <OperationalPage
      title="Import from Railway"
      description="Bring services, variables, domains, and database metadata from a Railway project into Vessl."
      scope="Import"
      statusLabel="Ready"
      statusTone="info"
      primaryAction={{ label: 'Start import', icon: Import }}
      metrics={[
        {
          label: 'Provider',
          value: 'Railway',
          detail: 'Project export',
          icon: TrainFront,
        },
        {
          label: 'Credentials',
          value: 'Required',
          detail: 'Token encrypted locally',
          icon: KeyRound,
        },
      ]}
    >
      <SettingsPanel
        title="Railway source"
        icon={UploadCloud}
        fields={[
          { label: 'Railway token', value: '' },
          { label: 'Project ID', value: '' },
          { label: 'Target Vessl project', value: '' },
        ]}
        toggles={[
          {
            label: 'Import environment variables',
            description: 'Copy non-empty variables into the target project after review.',
            checked: true,
          },
        ]}
        actionLabel="Validate source"
      />
    </OperationalPage>
  );
}
