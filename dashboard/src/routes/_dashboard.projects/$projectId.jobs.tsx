import { createFileRoute } from '@tanstack/react-router';
import { Clock3, Play, Plus } from 'lucide-react';
import { EmptyPanel } from '#/features/dashboard/empty-panel';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { ResourceTable } from '#/features/dashboard/resource-table';
import { StatusBadge } from '#/features/dashboard/status-badge';
import { useListJobs } from '#/hooks/useJobs';

export const Route = createFileRoute('/_dashboard/projects/$projectId/jobs')({
  component: ProjectJobsPage,
});

function ProjectJobsPage() {
  const { projectId } = Route.useParams();
  const jobs = useListJobs(projectId);
  const rows = jobs.data?.data ?? [];

  return (
    <OperationalPage
      title="Jobs"
      description="Schedule recurring commands, inspect run status, and trigger maintenance tasks for this project."
      scope={projectId}
      statusLabel={jobs.isError ? 'Jobs unavailable' : 'Scheduler ready'}
      statusTone={jobs.isError ? 'danger' : 'healthy'}
      primaryAction={{ label: 'New job', icon: Plus }}
      secondaryActions={[{ label: 'Run selected', icon: Play, variant: 'outline' }]}
      metrics={[
        {
          label: 'Jobs',
          value: jobs.isLoading ? undefined : rows.length.toString(),
          detail: 'Scheduled tasks',
          icon: Clock3,
        },
      ]}
    >
      {rows.length > 0 ? (
        <ResourceTable
          columns={[
            { key: 'name', label: 'Job', render: (job) => job.name },
            { key: 'schedule', label: 'Schedule', render: (job) => job.schedule },
            {
              key: 'status',
              label: 'Status',
              render: (job) => (
                <StatusBadge
                  label={job.status}
                  tone={job.status === 'failed' || job.status === 'error' ? 'danger' : 'healthy'}
                />
              ),
            },
            { key: 'lastRun', label: 'Last run', render: (job) => job.lastRunAt || 'Never' },
          ]}
          rows={rows}
          getRowKey={(job) => job.id}
        />
      ) : (
        <EmptyPanel
          icon={Clock3}
          title="No scheduled jobs"
          description="Create jobs for database cleanup, background processing, cron tasks, and other recurring project work."
          actionLabel="New job"
        />
      )}
    </OperationalPage>
  );
}
