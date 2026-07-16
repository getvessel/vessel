import { createFileRoute } from '@tanstack/react-router';
import { FolderKanban, GitBranch, Plus, ShieldCheck } from 'lucide-react';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { ProjectList } from '#/features/projects/project-list';
import { useListProjects } from '#/hooks/useProjects';

export const Route = createFileRoute('/_dashboard/projects')({
  component: ProjectsPage,
});

function ProjectsPage() {
  const projects = useListProjects();
  const rows = projects.data?.data.records ?? [];

  return (
    <OperationalPage
      title="Projects"
      description="Inspect workloads, environments, deploy targets, and infrastructure relationships across the instance."
      scope="Workspace"
      statusLabel={projects.isError ? 'Needs attention' : 'Operational'}
      statusTone={projects.isError ? 'danger' : 'healthy'}
      primaryAction={{ label: 'New project', icon: Plus }}
      secondaryActions={[{ label: 'Import repository', icon: GitBranch, variant: 'outline' }]}
      metrics={[
        {
          label: 'Projects',
          value: projects.isLoading ? undefined : rows.length.toString(),
          detail: 'Tracked workspaces',
          icon: FolderKanban,
        },
        {
          label: 'Policy',
          value: projects.isLoading ? undefined : 'Inherited',
          detail: 'Instance defaults applied',
          icon: ShieldCheck,
        },
      ]}
    >
      <ProjectList />
    </OperationalPage>
  );
}
