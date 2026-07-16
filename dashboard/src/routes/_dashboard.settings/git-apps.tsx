import { createFileRoute } from '@tanstack/react-router';
import { GitBranch, KeyRound, Plus } from 'lucide-react';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { ResourceTable } from '#/features/dashboard/resource-table';
import { StatusBadge } from '#/features/dashboard/status-badge';

export const Route = createFileRoute('/_dashboard/settings/git-apps')({
  component: GitAppsSettingsPage,
});

const gitProviderRows = [
  { id: 'github', provider: 'GitHub', status: 'Ready', appId: 'Not installed' },
  { id: 'gitlab', provider: 'GitLab', status: 'Optional', appId: 'Manual token' },
];

function GitAppsSettingsPage() {
  return (
    <OperationalPage
      title="Git apps"
      description="Connect repository providers for imports, deployment triggers, commit metadata, and preview environments."
      scope="Integrations"
      statusLabel="Providers available"
      statusTone="info"
      primaryAction={{ label: 'Add provider', icon: Plus }}
      metrics={[
        {
          label: 'Providers',
          value: gitProviderRows.length.toString(),
          detail: 'Supported integrations',
          icon: GitBranch,
          tone: 'info',
        },
        {
          label: 'Secrets',
          value: 'Encrypted',
          detail: 'Provider credentials',
          icon: KeyRound,
          tone: 'healthy',
        },
      ]}
    >
      <ResourceTable
        columns={[
          { key: 'provider', label: 'Provider', render: (provider) => provider.provider },
          { key: 'appId', label: 'Configuration', render: (provider) => provider.appId },
          {
            key: 'status',
            label: 'Status',
            render: (provider) => (
              <StatusBadge
                label={provider.status}
                tone={provider.status === 'Ready' ? 'healthy' : 'neutral'}
              />
            ),
          },
        ]}
        rows={gitProviderRows}
        getRowKey={(provider) => provider.id}
      />
    </OperationalPage>
  );
}
