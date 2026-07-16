import { createFileRoute } from '@tanstack/react-router';
import { FolderCog, Save, ShieldCheck } from 'lucide-react';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { SettingsPanel } from '#/features/dashboard/settings-panel';
import { useGetProject } from '#/hooks/useProjects';

export const Route = createFileRoute('/_dashboard/projects/$projectId/settings')({
  component: ProjectSettingsPage,
});

function ProjectSettingsPage() {
  const { projectId } = Route.useParams();
  const project = useGetProject(projectId);
  const data = project.data?.data;

  return (
    <OperationalPage
      title="Project settings"
      description="Update identity, defaults, access boundaries, and destructive controls for this project."
      scope={data?.name ?? projectId}
      statusLabel={project.isError ? 'Unavailable' : 'Editable'}
      statusTone={project.isError ? 'danger' : 'healthy'}
      primaryAction={{ label: 'Save settings', icon: Save }}
      metrics={[
        {
          label: 'Access',
          value: 'Team scoped',
          detail: 'Inherited permissions',
          icon: ShieldCheck,
        },
      ]}
    >
      <SettingsPanel
        title="Project identity"
        icon={FolderCog}
        fields={[
          { label: 'Name', value: data?.name ?? '' },
          { label: 'Description', value: data?.description ?? '', type: 'textarea' },
        ]}
        toggles={[
          {
            label: 'Protect production deployments',
            description: 'Require explicit confirmation before replacing production workloads.',
            checked: true,
          },
        ]}
      />
    </OperationalPage>
  );
}
