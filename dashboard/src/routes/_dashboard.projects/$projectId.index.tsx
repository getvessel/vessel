import { createFileRoute } from '@tanstack/react-router';
import { Boxes, Database, Plus, Settings } from 'lucide-react';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { ProjectCommandCenter } from '#/features/projects/project-command-center';
import { useListByProject } from '#/hooks/useApps';
import { useGetProject } from '#/hooks/useProjects';

export const Route = createFileRoute('/_dashboard/projects/$projectId/')({
  component: ProjectOverviewPage,
});

function ProjectOverviewPage() {
  const { projectId } = Route.useParams();
  const project = useGetProject(projectId);
  const services = useListByProject(projectId);
  const serviceRows = services.data?.data ?? [];

  return (
    <OperationalPage
      title={project.data?.data.name ?? 'Project'}
      description="Project workspace for services, databases, scheduled jobs, environment canvas, and configuration."
      scope={projectId}
      statusLabel={project.isError ? 'Unavailable' : 'Active'}
      statusTone={project.isError ? 'danger' : 'healthy'}
      primaryAction={{ label: 'Add service', icon: Plus }}
      secondaryActions={[{ label: 'Project settings', icon: Settings, variant: 'outline' }]}
      metrics={[
        {
          label: 'Services',
          value: services.isLoading ? undefined : serviceRows.length.toString(),
          detail: 'Application workloads',
          icon: Boxes,
        },
        {
          label: 'Databases',
          value: 'Linked',
          detail: 'Managed separately',
          icon: Database,
        },
      ]}
    >
      <ProjectCommandCenter projectId={projectId} services={serviceRows} />
    </OperationalPage>
  );
}
