import { Link } from '@tanstack/react-router';
import { Activity, Clock3, Database, GitBranch, Globe2, Plus, Server } from 'lucide-react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { EmptyPanel } from '#/features/dashboard/empty-panel';
import { EnvironmentSwitcher } from '#/features/dashboard/environment-switcher';
import { ResourceHealthStrip } from '#/features/dashboard/resource-health-strip';
import { ResourceTable } from '#/features/dashboard/resource-table';
import { StatusBadge } from '#/features/dashboard/status-badge';
import type { AppService } from '#/interfaces/deployment';
import { ProjectActivityFeed } from './project-activity-feed';
import { ProjectDependencyMap } from './project-dependency-map';

interface ProjectCommandCenterProps {
  projectId: string;
  services: AppService[];
}

export function ProjectCommandCenter({ projectId, services }: ProjectCommandCenterProps) {
  const runningServices = services.filter((service) => service.status === 'running').length;
  const queueAction = (label: string) => {
    toast.info(`${label} queued`, {
      description: `This workflow will start in ${projectId} when the endpoint is connected.`,
    });
  };

  return (
    <div className="grid gap-5">
      <div className="flex flex-col gap-3 rounded-lg border bg-card p-4 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <h2 className="font-semibold text-base">Production command center</h2>
          <p className="text-muted-foreground text-sm">
            Runtime health, services, dependencies, and recent operations.
          </p>
        </div>
        <EnvironmentSwitcher />
      </div>

      <ResourceHealthStrip
        items={[
          {
            label: 'Services',
            value: `${runningServices}/${services.length || 0} running`,
            detail: 'Application workloads',
            icon: Server,
            tone: runningServices === services.length && services.length > 0 ? 'healthy' : 'info',
          },
          {
            label: 'Deployments',
            value: services.length > 0 ? 'Live' : 'Idle',
            detail: 'Latest release channel',
            icon: Activity,
            tone: services.length > 0 ? 'healthy' : 'neutral',
          },
          {
            label: 'Domains',
            value: services.some((service) => service.domain) ? 'Configured' : 'Pending',
            detail: 'Public routing',
            icon: Globe2,
            tone: services.some((service) => service.domain) ? 'healthy' : 'warning',
          },
          {
            label: 'Database',
            value: 'Ready',
            detail: 'Managed separately',
            icon: Database,
            tone: 'healthy',
          },
        ]}
      />

      <div className="grid gap-5 xl:grid-cols-[minmax(0,1fr)_360px]">
        <div className="grid gap-5">
          {services.length > 0 ? (
            <ResourceTable
              caption="Project services"
              columns={[
                {
                  key: 'name',
                  label: 'Service',
                  render: (service) => (
                    <Link
                      to="/services/$serviceId"
                      params={{ serviceId: service.id }}
                      className="font-medium hover:text-primary"
                    >
                      {service.name}
                    </Link>
                  ),
                },
                {
                  key: 'status',
                  label: 'Status',
                  render: (service) => (
                    <StatusBadge
                      label={service.status}
                      tone={service.status === 'running' ? 'healthy' : 'warning'}
                    />
                  ),
                },
                { key: 'branch', label: 'Branch', render: (service) => service.branch || 'main' },
                {
                  key: 'domain',
                  label: 'Domain',
                  render: (service) => service.domain || 'Not assigned',
                },
              ]}
              rows={services}
              getRowKey={(service) => service.id}
            />
          ) : (
            <EmptyPanel
              icon={GitBranch}
              title="Create your first workload"
              description="Import a repository, deploy a container image, or start from a template to create the first service in this project."
              actionLabel="Add service"
            />
          )}
          <ProjectDependencyMap services={services} />
        </div>

        <div className="grid gap-5">
          <section className="rounded-lg border bg-card p-4">
            <h2 className="font-semibold text-base">Quick create</h2>
            <p className="text-muted-foreground text-sm">Start common project workflows.</p>
            <div className="mt-4 grid gap-2">
              {[
                { label: 'New service', icon: Plus },
                { label: 'Import repository', icon: GitBranch },
                { label: 'Create database', icon: Database },
                { label: 'Schedule job', icon: Clock3 },
              ].map((action) => (
                <Button
                  key={action.label}
                  type="button"
                  variant="outline"
                  className="justify-start"
                  onClick={() => queueAction(action.label)}
                >
                  <action.icon className="size-4" />
                  {action.label}
                </Button>
              ))}
            </div>
          </section>
          <ProjectActivityFeed />
        </div>
      </div>
    </div>
  );
}
