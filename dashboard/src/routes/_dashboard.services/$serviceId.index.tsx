import { createFileRoute } from '@tanstack/react-router';
import { Activity, Boxes, GitBranch, Globe2, RefreshCw, Server } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { ResourceTable } from '#/features/dashboard/resource-table';
import { StatusBadge } from '#/features/dashboard/status-badge';
import { ServiceMetricsPanel } from '#/features/services/service-metrics';
import { useGetApp } from '#/hooks/useApps';
import { useListByService } from '#/hooks/useDeployments';

export const Route = createFileRoute('/_dashboard/services/$serviceId/')({
  component: ServiceOverviewPage,
});

function ServiceOverviewPage() {
  const { serviceId } = Route.useParams();
  const service = useGetApp(serviceId);
  const deployments = useListByService(serviceId);
  const app = service.data?.data;
  const deploymentRows = deployments.data?.data.records ?? [];

  return (
    <OperationalPage
      title={app?.name ?? 'Service'}
      description="Runtime status, deployment activity, traffic endpoints, and resource pressure for this service."
      scope={app?.projectId ?? serviceId}
      statusLabel={app?.status ?? (service.isError ? 'Unavailable' : 'Loading')}
      statusTone={app?.status === 'running' ? 'healthy' : service.isError ? 'danger' : 'warning'}
      primaryAction={{ label: 'Redeploy', icon: RefreshCw }}
      secondaryActions={[{ label: 'Open domain', icon: Globe2, variant: 'outline' }]}
      metrics={[
        {
          label: 'Runtime',
          value: app?.runtimeMode,
          detail: app?.buildEngine ?? 'Build engine pending',
          icon: Server,
        },
        {
          label: 'Branch',
          value: app?.branch,
          detail: app?.repositoryUrl ?? 'Repository not configured',
          icon: GitBranch,
        },
        {
          label: 'Deployments',
          value: deployments.isLoading ? undefined : deploymentRows.length.toString(),
          detail: 'Recent releases',
          icon: Boxes,
        },
        {
          label: 'Port',
          value: app?.internalPort ? app.internalPort.toString() : undefined,
          detail: app?.healthCheckPath || 'Health check pending',
          icon: Activity,
        },
      ]}
    >
      <div className="grid gap-4 xl:grid-cols-[1fr_360px]">
        <div className="grid gap-4">
          <ServiceMetricsPanel serviceId={serviceId} />
          <ResourceTable
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
                render: (deployment) => deployment.commitHash || 'Manual trigger',
              },
              {
                key: 'updated',
                label: 'Updated',
                render: (deployment) => new Date(deployment.updatedAt).toLocaleString(),
              },
            ]}
            rows={deploymentRows}
            getRowKey={(deployment) => deployment.id}
          />
        </div>
        <Card className="shadow-none">
          <CardHeader>
            <CardTitle className="text-base">Service routing</CardTitle>
          </CardHeader>
          <CardContent className="grid gap-3 text-sm">
            <div className="flex justify-between gap-4">
              <span className="text-muted-foreground">Domain</span>
              <span className="truncate">{app?.domain || 'Not assigned'}</span>
            </div>
            <div className="flex justify-between gap-4">
              <span className="text-muted-foreground">Container</span>
              <span className="truncate">{app?.containerId || 'Pending'}</span>
            </div>
            <div className="flex justify-between gap-4">
              <span className="text-muted-foreground">Root</span>
              <span>{app?.rootDirectory || '/'}</span>
            </div>
          </CardContent>
        </Card>
      </div>
    </OperationalPage>
  );
}
