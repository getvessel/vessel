import { createFileRoute } from '@tanstack/react-router';
import { Plus, ShieldCheck, Users } from 'lucide-react';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { ResourceTable } from '#/features/dashboard/resource-table';
import { StatusBadge } from '#/features/dashboard/status-badge';

export const Route = createFileRoute('/_dashboard/settings/users')({
  component: UsersSettingsPage,
});

const userRows = [
  { id: 'owner', name: 'Instance owner', email: 'owner@local', role: 'Owner', status: 'Active' },
  { id: 'ops', name: 'Operations', email: 'ops@local', role: 'Admin', status: 'Invited' },
];

function UsersSettingsPage() {
  return (
    <OperationalPage
      title="Users"
      description="Invite operators, assign roles, and keep administrator access visible."
      scope="Access"
      statusLabel="Protected"
      statusTone="healthy"
      primaryAction={{ label: 'Invite user', icon: Plus }}
      metrics={[
        {
          label: 'Users',
          value: userRows.length.toString(),
          detail: 'Known operators',
          icon: Users,
        },
        {
          label: 'Roles',
          value: 'Scoped',
          detail: 'Least privilege',
          icon: ShieldCheck,
        },
      ]}
    >
      <ResourceTable
        columns={[
          { key: 'name', label: 'User', render: (user) => user.name },
          { key: 'email', label: 'Email', render: (user) => user.email },
          { key: 'role', label: 'Role', render: (user) => user.role },
          {
            key: 'status',
            label: 'Status',
            render: (user) => (
              <StatusBadge
                label={user.status}
                tone={user.status === 'Active' ? 'healthy' : 'warning'}
              />
            ),
          },
        ]}
        rows={userRows}
        getRowKey={(user) => user.id}
      />
    </OperationalPage>
  );
}
