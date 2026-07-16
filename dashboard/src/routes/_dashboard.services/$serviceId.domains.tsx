import { createFileRoute } from '@tanstack/react-router';
import { Globe2, LockKeyhole, Plus, Route as RouteIcon } from 'lucide-react';
import { EmptyPanel } from '#/features/dashboard/empty-panel';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { ResourceTable } from '#/features/dashboard/resource-table';
import { StatusBadge } from '#/features/dashboard/status-badge';
import { useGetApp } from '#/hooks/useApps';

export const Route = createFileRoute('/_dashboard/services/$serviceId/domains')({
  component: ServiceDomainsPage,
});

function ServiceDomainsPage() {
  const { serviceId } = Route.useParams();
  const service = useGetApp(serviceId);
  const app = service.data?.data;
  const rows = app?.domain ? [{ id: app.id, domain: app.domain, port: app.internalPort }] : [];

  return (
    <OperationalPage
      title="Domains"
      description="Route public traffic to this service, verify DNS, and track certificate readiness."
      scope={app?.name ?? serviceId}
      statusLabel={rows.length > 0 ? 'Route configured' : 'No domain'}
      statusTone={rows.length > 0 ? 'healthy' : 'warning'}
      primaryAction={{ label: 'Add domain', icon: Plus }}
      metrics={[
        {
          label: 'Domains',
          value: service.isLoading ? undefined : rows.length.toString(),
          detail: 'Public routes',
          icon: Globe2,
        },
        {
          label: 'TLS',
          value: rows.length > 0 ? 'Pending' : 'Inactive',
          detail: 'Certificate status',
          icon: LockKeyhole,
        },
      ]}
    >
      {rows.length > 0 ? (
        <ResourceTable
          columns={[
            { key: 'domain', label: 'Domain', render: (domain) => domain.domain },
            { key: 'target', label: 'Target', render: (domain) => `:${domain.port}` },
            {
              key: 'status',
              label: 'Certificate',
              render: () => <StatusBadge label="Provisioning" tone="warning" />,
            },
          ]}
          rows={rows}
          getRowKey={(domain) => domain.id}
        />
      ) : (
        <EmptyPanel
          icon={RouteIcon}
          title="No domains attached"
          description="Attach a custom domain to expose this service and let Vessl handle routing and certificate lifecycle."
          actionLabel="Add domain"
        />
      )}
    </OperationalPage>
  );
}
