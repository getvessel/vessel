import { createFileRoute } from '@tanstack/react-router';
import { KeyRound, Plus, ShieldCheck } from 'lucide-react';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { ResourceTable } from '#/features/dashboard/resource-table';
import { StatusBadge } from '#/features/dashboard/status-badge';

export const Route = createFileRoute('/_dashboard/settings/oauth')({
  component: OAuthSettingsPage,
});

const providerRows = [
  { id: 'github', provider: 'GitHub', status: 'Optional', callback: '/api/auth/github/callback' },
  { id: 'google', provider: 'Google', status: 'Optional', callback: '/api/auth/google/callback' },
];

function OAuthSettingsPage() {
  return (
    <OperationalPage
      title="OAuth"
      description="Configure external identity providers and keep authentication policy visible to administrators."
      scope="Identity"
      statusLabel="Password login active"
      statusTone="healthy"
      primaryAction={{ label: 'Add provider', icon: Plus }}
      metrics={[
        {
          label: 'Providers',
          value: providerRows.length.toString(),
          detail: 'Available connectors',
          icon: KeyRound,
        },
        {
          label: 'Policy',
          value: 'Admin managed',
          detail: 'Login controls',
          icon: ShieldCheck,
        },
      ]}
    >
      <ResourceTable
        columns={[
          { key: 'provider', label: 'Provider', render: (provider) => provider.provider },
          { key: 'callback', label: 'Callback', render: (provider) => provider.callback },
          {
            key: 'status',
            label: 'Status',
            render: (provider) => <StatusBadge label={provider.status} tone="neutral" />,
          },
        ]}
        rows={providerRows}
        getRowKey={(provider) => provider.id}
      />
    </OperationalPage>
  );
}
