import { createFileRoute } from '@tanstack/react-router';
import { GitCommit, RefreshCw, RotateCcw, Timer } from 'lucide-react';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { ServiceReleaseCenter } from '#/features/services/service-release-center';
import { useListByService } from '#/hooks/useDeployments';

export const Route = createFileRoute('/_dashboard/services/$serviceId/deployments')({
  component: ServiceDeploymentsPage,
});

function ServiceDeploymentsPage() {
  const { serviceId } = Route.useParams();
  const deployments = useListByService(serviceId);
  const rows = deployments.data?.data.records ?? [];

  return (
    <OperationalPage
      title="Deployments"
      description="Release history with status, commit context, duration, rollback, and redeploy actions."
      scope={serviceId}
      statusLabel={deployments.isError ? 'History unavailable' : 'Tracking releases'}
      statusTone={deployments.isError ? 'danger' : 'healthy'}
      primaryAction={{ label: 'Redeploy latest', icon: RefreshCw }}
      secondaryActions={[{ label: 'Rollback', icon: RotateCcw, variant: 'outline' }]}
      metrics={[
        {
          label: 'Deployments',
          value: deployments.isLoading ? undefined : rows.length.toString(),
          detail: 'Recent records',
          icon: GitCommit,
        },
        {
          label: 'Latest',
          value: rows[0]?.status,
          detail: rows[0]?.updatedAt ? new Date(rows[0].updatedAt).toLocaleString() : 'No releases',
          icon: Timer,
        },
      ]}
    >
      <ServiceReleaseCenter deployments={rows} />
    </OperationalPage>
  );
}
