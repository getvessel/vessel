import { createFileRoute } from '@tanstack/react-router';
import { PlugZap, Terminal } from 'lucide-react';
import { LogPanel } from '#/features/dashboard/log-panel';
import { OperationalPage } from '#/features/dashboard/operational-page';

export const Route = createFileRoute('/_dashboard/services/$serviceId/terminal')({
  component: ServiceTerminalPage,
});

function ServiceTerminalPage() {
  const { serviceId } = Route.useParams();

  return (
    <OperationalPage
      title="Terminal"
      description="Open an interactive service shell for diagnostics, filesystem inspection, and one-off commands."
      scope={serviceId}
      statusLabel="Awaiting session"
      statusTone="warning"
      primaryAction={{ label: 'Connect shell', icon: PlugZap }}
      metrics={[
        {
          label: 'Session',
          value: 'Detached',
          detail: 'No shell attached',
          icon: Terminal,
        },
      ]}
    >
      <LogPanel
        title="Shell preview"
        lines={[
          '$ vessl exec service --interactive',
          'Waiting for daemon connection...',
          'No terminal session is attached yet.',
        ]}
      />
    </OperationalPage>
  );
}
