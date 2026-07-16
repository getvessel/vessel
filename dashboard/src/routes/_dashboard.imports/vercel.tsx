import { createFileRoute } from '@tanstack/react-router';
import { Import, KeyRound, Triangle, UploadCloud } from 'lucide-react';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { SettingsPanel } from '#/features/dashboard/settings-panel';

export const Route = createFileRoute('/_dashboard/imports/vercel')({
  component: VercelImportPage,
});

function VercelImportPage() {
  return (
    <OperationalPage
      title="Import from Vercel"
      description="Convert Vercel project settings into Vessl services, variables, domains, and build commands."
      scope="Import"
      statusLabel="Ready"
      statusTone="info"
      primaryAction={{ label: 'Start import', icon: Import }}
      metrics={[
        {
          label: 'Provider',
          value: 'Vercel',
          detail: 'Project export',
          icon: Triangle,
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
        title="Vercel source"
        icon={UploadCloud}
        fields={[
          { label: 'Vercel token', value: '' },
          { label: 'Team or user slug', value: '' },
          { label: 'Project name', value: '' },
        ]}
        toggles={[
          {
            label: 'Preserve production domain',
            description: 'Create the imported domain as pending until DNS points at this instance.',
            checked: true,
          },
        ]}
        actionLabel="Validate source"
      />
    </OperationalPage>
  );
}
