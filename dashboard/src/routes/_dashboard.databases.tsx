import { createFileRoute } from '@tanstack/react-router';
import { Database, Plus, ShieldCheck, Waypoints } from 'lucide-react';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { DatabaseList } from '#/features/databases/database-list';
import { useListDatabases } from '#/hooks/useDatabases';

export const Route = createFileRoute('/_dashboard/databases')({
  component: DatabasesPage,
});

function DatabasesPage() {
  const databases = useListDatabases();
  const rows = databases.data?.data ?? [];
  const runningCount = rows.filter((database) => database.status === 'running').length;

  return (
    <OperationalPage
      title="Databases"
      description="Manage persistent services, connection details, backup posture, and SQL access from one inventory."
      scope="Storage"
      statusLabel={databases.isError ? 'API error' : 'Inventory live'}
      statusTone={databases.isError ? 'danger' : 'healthy'}
      primaryAction={{ label: 'New database', icon: Plus }}
      secondaryActions={[{ label: 'Import data', icon: Waypoints, variant: 'outline' }]}
      metrics={[
        {
          label: 'Total',
          value: databases.isLoading ? undefined : rows.length.toString(),
          detail: 'Managed databases',
          icon: Database,
        },
        {
          label: 'Running',
          value: databases.isLoading ? undefined : runningCount.toString(),
          detail: 'Accepting traffic',
          icon: ShieldCheck,
        },
      ]}
    >
      <DatabaseList />
    </OperationalPage>
  );
}
