import { createFileRoute } from '@tanstack/react-router';
import { Box, GitBranch, Hammer, PackageCheck } from 'lucide-react';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { SettingsPanel } from '#/features/dashboard/settings-panel';
import { useGetApp } from '#/hooks/useApps';

export const Route = createFileRoute('/_dashboard/services/$serviceId/build')({
  component: ServiceBuildPage,
});

function ServiceBuildPage() {
  const { serviceId } = Route.useParams();
  const service = useGetApp(serviceId);
  const app = service.data?.data;

  return (
    <OperationalPage
      title="Build settings"
      description="Control source, build engine, commands, output paths, and runtime startup behavior for this service."
      scope={app?.name ?? serviceId}
      statusLabel={app?.buildEngine ?? 'Configured'}
      statusTone={service.isError ? 'danger' : 'info'}
      primaryAction={{ label: 'Save build', icon: PackageCheck }}
      metrics={[
        {
          label: 'Engine',
          value: app?.buildEngine,
          detail: 'Build strategy',
          icon: Hammer,
        },
        {
          label: 'Branch',
          value: app?.branch,
          detail: app?.repositoryUrl || 'Repository source',
          icon: GitBranch,
        },
      ]}
    >
      <SettingsPanel
        title="Build pipeline"
        icon={Box}
        fields={[
          { label: 'Install command', value: app?.installCommand ?? '' },
          { label: 'Build command', value: app?.buildCommand ?? '' },
          { label: 'Start command', value: app?.startCommand ?? '' },
          { label: 'Dockerfile path', value: app?.dockerfilePath ?? 'Dockerfile' },
          { label: 'Static output directory', value: app?.staticOutput ?? 'dist' },
        ]}
        toggles={[
          {
            label: 'Use repository root',
            description:
              'Build from the configured root directory unless an environment overrides it.',
            checked: true,
          },
        ]}
      />
    </OperationalPage>
  );
}
