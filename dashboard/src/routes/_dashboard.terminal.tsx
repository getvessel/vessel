import { createFileRoute } from '@tanstack/react-router';
import { PlugZap, Terminal } from 'lucide-react';
import { LogPanel } from '#/features/dashboard/log-panel';
import { OperationalPage } from '#/features/dashboard/operational-page';

export const Route = createFileRoute('/_dashboard/terminal')({
  component: TerminalPage,
});

function TerminalPage() {
  return (
    <OperationalPage
      title="Instance terminal"
      description="Connect to the daemon host for diagnostics, emergency maintenance, and operational commands."
      scope="Operations"
      statusLabel="Detached"
      statusTone="warning"
      primaryAction={{ label: 'Connect', icon: PlugZap }}
      metrics={[
        {
          label: 'Shell',
          value: 'Host',
          detail: 'Interactive session',
          icon: Terminal,
        },
      ]}
    >
      <LogPanel
        title="Daemon terminal"
        lines={[
          '$ vessl daemon shell',
          'Interactive host sessions are created on demand.',
          'Connect to begin a controlled maintenance session.',
        ]}
      />
    </OperationalPage>
  );
}
