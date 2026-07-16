import { createFileRoute } from '@tanstack/react-router';
import { FileClock, Filter, ShieldCheck } from 'lucide-react';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { ResourceTable } from '#/features/dashboard/resource-table';
import { StatusBadge } from '#/features/dashboard/status-badge';

export const Route = createFileRoute('/_dashboard/audit-logs')({
  component: AuditLogsPage,
});

const auditRows = [
  {
    id: 'login',
    event: 'User signed in',
    actor: 'operator@localhost',
    target: 'dashboard',
    severity: 'Info',
  },
  {
    id: 'settings',
    event: 'Settings viewed',
    actor: 'operator@localhost',
    target: 'instance',
    severity: 'Info',
  },
];

function AuditLogsPage() {
  return (
    <OperationalPage
      title="Audit logs"
      description="Review administrative actions, authentication events, and sensitive changes across the instance."
      scope="Security"
      statusLabel="Recording"
      statusTone="healthy"
      secondaryActions={[{ label: 'Filter', icon: Filter, variant: 'outline' }]}
      metrics={[
        {
          label: 'Events',
          value: auditRows.length.toString(),
          detail: 'Recent activity',
          icon: FileClock,
        },
        {
          label: 'Integrity',
          value: 'Append-only',
          detail: 'Tamper-resistant log',
          icon: ShieldCheck,
        },
      ]}
    >
      <ResourceTable
        columns={[
          { key: 'event', label: 'Event', render: (event) => event.event },
          { key: 'actor', label: 'Actor', render: (event) => event.actor },
          { key: 'target', label: 'Target', render: (event) => event.target },
          {
            key: 'severity',
            label: 'Severity',
            render: (event) => <StatusBadge label={event.severity} tone="info" />,
          },
        ]}
        rows={auditRows}
        getRowKey={(event) => event.id}
      />
    </OperationalPage>
  );
}
