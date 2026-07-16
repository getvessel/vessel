import { createFileRoute } from '@tanstack/react-router';
import { KeyRound, Save, ShieldCheck, UserRound } from 'lucide-react';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { SettingsPanel } from '#/features/dashboard/settings-panel';

export const Route = createFileRoute('/_dashboard/profile/')({
  component: ProfilePage,
});

function ProfilePage() {
  return (
    <OperationalPage
      title="Profile"
      description="Manage personal identity, access tokens, session security, and notification preferences."
      scope="Account"
      statusLabel="Signed in"
      statusTone="healthy"
      primaryAction={{ label: 'Save profile', icon: Save }}
      metrics={[
        {
          label: 'Security',
          value: 'Password',
          detail: 'Primary login method',
          icon: ShieldCheck,
        },
        {
          label: 'Tokens',
          value: '0',
          detail: 'Personal API tokens',
          icon: KeyRound,
        },
      ]}
    >
      <SettingsPanel
        title="Account details"
        icon={UserRound}
        fields={[
          { label: 'Display name', value: 'Vessl operator' },
          { label: 'Email', value: 'operator@localhost' },
        ]}
        toggles={[
          {
            label: 'Require two-factor authentication',
            description: 'Add a second verification step for sensitive account actions.',
            checked: false,
          },
        ]}
      />
    </OperationalPage>
  );
}
