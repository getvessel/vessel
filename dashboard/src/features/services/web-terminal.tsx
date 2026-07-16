import { PlugZap, Terminal } from 'lucide-react';
import { LogPanel } from '#/features/dashboard/log-panel';
import { OperationalPage } from '#/features/dashboard/operational-page';

export const TerminalWindow = () => (
  <OperationalPage
    title="Terminal"
    description="Attach an interactive shell to inspect a running workload."
    scope="Service"
    statusLabel="Detached"
    statusTone="warning"
    primaryAction={{ label: 'Connect shell', icon: PlugZap }}
    metrics={[
      {
        label: 'Session',
        value: 'Idle',
        detail: 'No shell attached',
        icon: Terminal,
      },
    ]}
  >
    <LogPanel title="Shell" lines={['$ vessl exec service', 'Waiting for a terminal session...']} />
  </OperationalPage>
);
