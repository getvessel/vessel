import { Link } from '@tanstack/react-router';
import { ExternalLink, FolderKanban, GitBranch } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { Skeleton } from '#/components/ui/skeleton';
import { EmptyPanel } from '#/features/dashboard/empty-panel';
import { type ResourceColumn, ResourceTable } from '#/features/dashboard/resource-table';
import { StatusBadge } from '#/features/dashboard/status-badge';
import { useListProjects } from '#/hooks/useProjects';
import type { ProjectConfig } from '#/interfaces/project';

const projectColumns: ResourceColumn<ProjectConfig>[] = [
  {
    key: 'name',
    label: 'Project',
    render: (project) => (
      <div className="space-y-1">
        <Link
          to="/projects/$projectId"
          params={{ projectId: project.id }}
          className="font-medium hover:underline"
        >
          {project.name}
        </Link>
        <p className="max-w-xl truncate text-muted-foreground text-xs">
          {project.description || 'No description provided'}
        </p>
      </div>
    ),
  },
  {
    key: 'status',
    label: 'Health',
    render: () => <StatusBadge label="Configured" tone="healthy" />,
  },
  {
    key: 'updatedAt',
    label: 'Updated',
    render: (project) => (
      <span className="text-muted-foreground text-sm">
        {new Date(project.updatedAt).toLocaleString()}
      </span>
    ),
  },
  {
    key: 'actions',
    label: '',
    render: (project) => (
      <Button asChild size="sm" variant="ghost">
        <Link to="/projects/$projectId" params={{ projectId: project.id }}>
          <ExternalLink className="size-4" />
          Open
        </Link>
      </Button>
    ),
  },
];

export const ProjectList = () => {
  const projects = useListProjects();

  if (projects.isLoading) {
    return (
      <div className="grid gap-2 rounded-lg border p-4">
        {Array.from({ length: 5 }).map((_, index) => (
          <Skeleton key={index.toString()} className="h-12 w-full" />
        ))}
      </div>
    );
  }

  if (projects.isError) {
    return (
      <EmptyPanel
        icon={GitBranch}
        title="Projects could not be loaded"
        description="The dashboard could not reach the daemon project API. Keep the daemon running and refresh this page."
        actionLabel="Retry"
      />
    );
  }

  const rows = projects.data?.data.records ?? [];

  if (rows.length === 0) {
    return (
      <EmptyPanel
        icon={FolderKanban}
        title="Create your first project"
        description="Projects group services, databases, jobs, domains, and environments so each workload has a clear operational boundary."
        actionLabel="New project"
      />
    );
  }

  return <ResourceTable columns={projectColumns} rows={rows} getRowKey={(project) => project.id} />;
};
