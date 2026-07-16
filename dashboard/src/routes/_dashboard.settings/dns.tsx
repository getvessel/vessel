import { createFileRoute } from '@tanstack/react-router';
import { Globe2, Save, ShieldCheck } from 'lucide-react';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { SettingsPanel } from '#/features/dashboard/settings-panel';

export const Route = createFileRoute('/_dashboard/settings/dns')({
  component: DnsSettingsPage,
});

function DnsSettingsPage() {
  return (
    <OperationalPage
      title="DNS"
      description="Configure public hostnames, wildcard routing, certificate email, and ingress defaults."
      scope="Networking"
      statusLabel="Routable"
      statusTone="healthy"
      primaryAction={{ label: 'Save DNS', icon: Save }}
      metrics={[
        {
          label: 'Wildcard',
          value: 'Enabled',
          detail: 'Default app routing',
          icon: Globe2,
        },
        {
          label: 'TLS',
          value: 'Automatic',
          detail: 'Certificate lifecycle',
          icon: ShieldCheck,
        },
      ]}
    >
      <SettingsPanel
        title="Routing defaults"
        icon={Globe2}
        fields={[
          { label: 'Root domain', value: 'apps.localhost' },
          { label: 'Certificate email', value: 'admin@example.com' },
          { label: 'Ingress network', value: 'vessl-public' },
        ]}
        toggles={[
          {
            label: 'Enable wildcard routing',
            description: 'Automatically route service subdomains through the instance gateway.',
            checked: true,
          },
        ]}
      />
    </OperationalPage>
  );
}
