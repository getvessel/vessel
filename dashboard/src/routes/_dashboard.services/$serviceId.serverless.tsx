import { createFileRoute } from '@tanstack/react-router';
import { Code2, Cpu, Plus, TimerReset } from 'lucide-react';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { SettingsPanel } from '#/features/dashboard/settings-panel';
import { useGetApp } from '#/hooks/useApps';

export const Route = createFileRoute('/_dashboard/services/$serviceId/serverless')({
  component: ServerlessPage,
});

function ServerlessPage() {
  const { serviceId } = Route.useParams();
  const service = useGetApp(serviceId);
  const app = service.data?.data;

  return (
    <OperationalPage
      title="Serverless"
      description="Configure function runtime behavior, timeout, scaling boundaries, and source entrypoints."
      scope={app?.name ?? serviceId}
      statusLabel={app?.buildEngine === 'serverless' ? 'Enabled' : 'Optional'}
      statusTone={app?.buildEngine === 'serverless' ? 'healthy' : 'neutral'}
      primaryAction={{ label: 'Add function', icon: Plus }}
      metrics={[
        {
          label: 'Runtime',
          value: app?.buildEngine === 'serverless' ? app.runtimeMode : 'Disabled',
          detail: 'Execution mode',
          icon: Code2,
        },
        {
          label: 'Timeout',
          value: '30s',
          detail: 'Default request limit',
          icon: TimerReset,
        },
      ]}
    >
      <SettingsPanel
        title="Function defaults"
        icon={Cpu}
        fields={[
          { label: 'Entrypoint', value: 'handler.ts' },
          { label: 'Memory limit', value: '512 MB' },
          { label: 'Timeout', value: '30 seconds' },
        ]}
        toggles={[
          {
            label: 'Scale to zero',
            description: 'Pause idle functions and resume them on the next request.',
            checked: true,
          },
        ]}
      />
    </OperationalPage>
  );
}
