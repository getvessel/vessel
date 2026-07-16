import { createFileRoute } from '@tanstack/react-router';
import { Plus, ShieldCheck, UsersRound } from 'lucide-react';
import { EmptyPanel } from '#/features/dashboard/empty-panel';
import { OperationalPage } from '#/features/dashboard/operational-page';

export const Route = createFileRoute('/_dashboard/teams')({
  component: TeamsPage,
});

function TeamsPage() {
  return (
    <OperationalPage
      title="Teams"
      description="Organize users, projects, and permissions around operational ownership."
      scope="Access"
      statusLabel="Single team"
      statusTone="info"
      primaryAction={{ label: 'Create team', icon: Plus }}
      metrics={[
        {
          label: 'Teams',
          value: '1',
          detail: 'Default workspace',
          icon: UsersRound,
        },
        {
          label: 'Permissions',
          value: 'Scoped',
          detail: 'Project access boundaries',
          icon: ShieldCheck,
        },
      ]}
    >
      <EmptyPanel
        icon={UsersRound}
        title="Team management is ready"
        description="Create teams when multiple groups need separate project ownership, billing boundaries, or access policies."
        actionLabel="Create team"
      />
    </OperationalPage>
  );
}
