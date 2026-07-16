import { createFileRoute } from '@tanstack/react-router';
import { DownloadCloud, RefreshCw, Rocket } from 'lucide-react';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { ResourceTable } from '#/features/dashboard/resource-table';
import { StatusBadge } from '#/features/dashboard/status-badge';
import { useGetUpdateStatus } from '#/hooks/useSettings';

export const Route = createFileRoute('/_dashboard/settings/updates')({
  component: UpdatesSettingsPage,
});

function UpdatesSettingsPage() {
  const update = useGetUpdateStatus();
  const rows = [
    {
      id: 'control-plane',
      component: 'Control plane',
      current: 'local',
      latest: update.isError ? 'Unavailable' : 'Check required',
      status: update.isError ? 'Error' : 'Current',
    },
  ];

  return (
    <OperationalPage
      title="Updates"
      description="Check release availability, stage control plane upgrades, and track update rollout state."
      scope="Maintenance"
      statusLabel={update.isError ? 'Check failed' : 'Ready'}
      statusTone={update.isError ? 'danger' : 'healthy'}
      primaryAction={{ label: 'Check updates', icon: RefreshCw }}
      secondaryActions={[{ label: 'Deploy update', icon: Rocket, variant: 'outline' }]}
      metrics={[
        {
          label: 'Channel',
          value: 'Stable',
          detail: 'Release track',
          icon: DownloadCloud,
        },
      ]}
    >
      <ResourceTable
        columns={[
          { key: 'component', label: 'Component', render: (component) => component.component },
          { key: 'current', label: 'Current', render: (component) => component.current },
          { key: 'latest', label: 'Latest', render: (component) => component.latest },
          {
            key: 'status',
            label: 'Status',
            render: (component) => (
              <StatusBadge
                label={component.status}
                tone={component.status === 'Error' ? 'danger' : 'healthy'}
              />
            ),
          },
        ]}
        rows={rows}
        getRowKey={(component) => component.id}
      />
    </OperationalPage>
  );
}
