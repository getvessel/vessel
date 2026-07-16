import { Clock3, GitBranch, GitCommit, Rocket } from 'lucide-react';
import { EmptyPanel } from '#/features/dashboard/empty-panel';
import { EnvironmentSwitcher } from '#/features/dashboard/environment-switcher';
import { ResourceHealthStrip } from '#/features/dashboard/resource-health-strip';
import { ResourceTable } from '#/features/dashboard/resource-table';
import { StatusBadge } from '#/features/dashboard/status-badge';
import type { Deployment } from '#/interfaces/deployment';
import { DeployPreviewCard } from './deploy-preview-card';
import { DeploymentLogConsole } from './deployment-log-console';
import { DeploymentTimeline } from './deployment-timeline';

interface ServiceReleaseCenterProps {
  deployments: Deployment[];
}

export function ServiceReleaseCenter({ deployments }: ServiceReleaseCenterProps) {
  const latest = deployments[0];
  const successfulDeployments = deployments.filter((deployment) =>
    ['READY', 'ACTIVE', 'SUCCESS'].includes(deployment.status)
  ).length;

  if (!latest) {
    return (
      <EmptyPanel
        icon={Rocket}
        title="No deployments yet"
        description="Trigger a deployment to build the service and start tracking commit history, logs, duration, and rollback points."
        actionLabel="Deploy service"
      />
    );
  }

  return (
    <div className="grid gap-5">
      <div className="flex flex-col gap-3 rounded-lg border bg-card p-4 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <h2 className="font-semibold text-base">Release command center</h2>
          <p className="text-muted-foreground text-sm">
            Deployment previews, progress, logs, and rollback context.
          </p>
        </div>
        <EnvironmentSwitcher />
      </div>

      <ResourceHealthStrip
        items={[
          {
            label: 'Latest deploy',
            value: latest.status,
            detail: new Date(latest.updatedAt).toLocaleString(),
            icon: Rocket,
            tone: latest.status === 'FAILED' ? 'danger' : 'healthy',
          },
          {
            label: 'Successful',
            value: successfulDeployments.toString(),
            detail: 'Rollback candidates',
            icon: GitCommit,
            tone: 'healthy',
          },
          {
            label: 'Branch',
            value: latest.branch || 'main',
            detail: latest.commitHash?.slice(0, 7) || 'Manual trigger',
            icon: GitBranch,
            tone: 'info',
          },
          {
            label: 'Duration',
            value: latest.finishedAt ? '<1m' : 'Running',
            detail: 'Build and release',
            icon: Clock3,
            tone: latest.finishedAt ? 'neutral' : 'warning',
          },
        ]}
      />

      <DeployPreviewCard deployment={latest} />

      <div className="grid gap-5 xl:grid-cols-[360px_minmax(0,1fr)]">
        <DeploymentTimeline status={latest.status} />
        <DeploymentLogConsole />
      </div>

      <ResourceTable
        caption="Deployment history"
        columns={[
          {
            key: 'status',
            label: 'Status',
            render: (deployment) => (
              <StatusBadge
                label={deployment.status}
                tone={deployment.status === 'FAILED' ? 'danger' : 'healthy'}
              />
            ),
          },
          {
            key: 'commit',
            label: 'Commit',
            render: (deployment) =>
              deployment.commitMessage || deployment.commitHash || 'Manual deployment',
          },
          {
            key: 'trigger',
            label: 'Trigger',
            render: (deployment) => deployment.trigger || deployment.branch || 'dashboard',
          },
          {
            key: 'updated',
            label: 'Updated',
            render: (deployment) => new Date(deployment.updatedAt).toLocaleString(),
          },
        ]}
        rows={deployments}
        getRowKey={(deployment) => deployment.id}
      />
    </div>
  );
}
