import { createFileRoute } from '@tanstack/react-router';
import { Database, HardDrive, Network, ShieldCheck } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { StatusBadge } from '#/features/dashboard/status-badge';
import { useGetDatabase } from '#/hooks/useDatabases';

export const Route = createFileRoute('/_dashboard/databases/$dbId/')({
  component: DatabaseOverviewPage,
});

function DatabaseOverviewPage() {
  const { dbId } = Route.useParams();
  const database = useGetDatabase(dbId);
  const db = database.data?.data;

  return (
    <OperationalPage
      title={db?.name ?? 'Database'}
      description="Connection details, runtime health, replication posture, and storage location for this managed database."
      scope={db?.engine ?? dbId}
      statusLabel={db?.status ?? (database.isError ? 'Unavailable' : 'Loading')}
      statusTone={db?.status === 'running' ? 'healthy' : database.isError ? 'danger' : 'warning'}
      metrics={[
        {
          label: 'Engine',
          value: db?.engine,
          detail: db?.version ?? 'Version pending',
          icon: Database,
        },
        {
          label: 'Port',
          value: db?.port ? db.port.toString() : undefined,
          detail: 'Internal listener',
          icon: Network,
        },
        {
          label: 'Replication',
          value: db ? (db.logicalReplication ? 'Enabled' : 'Disabled') : undefined,
          detail: 'Logical replication',
          icon: ShieldCheck,
        },
        {
          label: 'Storage',
          value: db?.volumePath ? 'Mounted' : undefined,
          detail: db?.volumePath ?? 'Volume pending',
          icon: HardDrive,
        },
      ]}
    >
      <div className="grid gap-4 lg:grid-cols-2">
        <Card className="shadow-none">
          <CardHeader>
            <CardTitle className="text-base">Connection</CardTitle>
          </CardHeader>
          <CardContent className="grid gap-3 text-sm">
            <div className="flex justify-between gap-4">
              <span className="text-muted-foreground">Internal DNS</span>
              <span className="truncate font-mono">{db?.internalDns || 'Pending'}</span>
            </div>
            <div className="flex justify-between gap-4">
              <span className="text-muted-foreground">External DNS</span>
              <span className="truncate font-mono">{db?.externalDns || 'Disabled'}</span>
            </div>
            <div className="flex justify-between gap-4">
              <span className="text-muted-foreground">Database</span>
              <span className="font-mono">{db?.databaseName || 'Pending'}</span>
            </div>
          </CardContent>
        </Card>
        <Card className="shadow-none">
          <CardHeader>
            <CardTitle className="text-base">Health checks</CardTitle>
          </CardHeader>
          <CardContent className="grid gap-3 text-sm">
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">Container</span>
              <StatusBadge label={db?.containerId ? 'Attached' : 'Pending'} tone="info" />
            </div>
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">Readiness</span>
              <StatusBadge
                label={db?.status === 'running' ? 'Passing' : 'Waiting'}
                tone={db?.status === 'running' ? 'healthy' : 'warning'}
              />
            </div>
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">Pressure</span>
              <StatusBadge label="Nominal" tone="healthy" />
            </div>
          </CardContent>
        </Card>
      </div>
    </OperationalPage>
  );
}
