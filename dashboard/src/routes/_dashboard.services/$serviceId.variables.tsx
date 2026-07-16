import { createFileRoute } from '@tanstack/react-router';
import { KeyRound, Plus, ShieldCheck, Variable } from 'lucide-react';
import { EmptyPanel } from '#/features/dashboard/empty-panel';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { ResourceTable } from '#/features/dashboard/resource-table';
import { StatusBadge } from '#/features/dashboard/status-badge';
import { useGetVars } from '#/hooks/useProjects';

export const Route = createFileRoute('/_dashboard/services/$serviceId/variables')({
  component: ServiceVariablesPage,
});

function ServiceVariablesPage() {
  const { serviceId } = Route.useParams();
  const variables = useGetVars(serviceId);
  const rows = Object.entries(variables.data?.data ?? {}).map(([key, value]) => ({ key, value }));

  return (
    <OperationalPage
      title="Variables"
      description="Manage runtime configuration and secrets that are injected into deployments."
      scope={serviceId}
      statusLabel={variables.isError ? 'Unable to sync' : 'Synced'}
      statusTone={variables.isError ? 'danger' : 'healthy'}
      primaryAction={{ label: 'Add variable', icon: Plus }}
      metrics={[
        {
          label: 'Variables',
          value: variables.isLoading ? undefined : rows.length.toString(),
          detail: 'Runtime keys',
          icon: Variable,
        },
        {
          label: 'Secrets',
          value: variables.isLoading ? undefined : rows.length.toString(),
          detail: 'Values masked at rest',
          icon: ShieldCheck,
        },
      ]}
    >
      {rows.length > 0 ? (
        <ResourceTable
          columns={[
            {
              key: 'key',
              label: 'Key',
              render: (variable) => <span className="font-mono text-sm">{variable.key}</span>,
            },
            {
              key: 'value',
              label: 'Value',
              render: () => (
                <span className="font-mono text-muted-foreground text-sm">••••••••••••</span>
              ),
            },
            {
              key: 'state',
              label: 'State',
              render: () => <StatusBadge label="Encrypted" tone="healthy" />,
            },
          ]}
          rows={rows}
          getRowKey={(variable) => variable.key}
        />
      ) : (
        <EmptyPanel
          icon={KeyRound}
          title="No variables configured"
          description="Add environment variables to keep credentials, feature flags, and service endpoints out of source control."
          actionLabel="Add variable"
        />
      )}
    </OperationalPage>
  );
}
