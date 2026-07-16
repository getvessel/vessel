import { Link } from '@tanstack/react-router';
import { Database, ExternalLink, ServerCrash } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { Skeleton } from '#/components/ui/skeleton';
import { EmptyPanel } from '#/features/dashboard/empty-panel';
import { type ResourceColumn, ResourceTable } from '#/features/dashboard/resource-table';
import { StatusBadge } from '#/features/dashboard/status-badge';
import { useListDatabases } from '#/hooks/useDatabases';
import type { Database as DatabaseResource } from '#/interfaces/database';

const databaseColumns: ResourceColumn<DatabaseResource>[] = [
  {
    key: 'name',
    label: 'Database',
    render: (database) => (
      <div className="space-y-1">
        <Link
          to="/databases/$dbId"
          params={{ dbId: database.id }}
          className="font-medium hover:underline"
        >
          {database.name}
        </Link>
        <p className="text-muted-foreground text-xs">
          {database.engine} {database.version}
        </p>
      </div>
    ),
  },
  {
    key: 'status',
    label: 'Status',
    render: (database) => (
      <StatusBadge
        label={database.status}
        tone={
          database.status === 'running'
            ? 'healthy'
            : database.status === 'error'
              ? 'danger'
              : 'warning'
        }
      />
    ),
  },
  {
    key: 'network',
    label: 'Network',
    render: (database) => (
      <div className="space-y-1 text-sm">
        <p>{database.internalDns || 'Internal DNS pending'}</p>
        <p className="text-muted-foreground text-xs">Port {database.port}</p>
      </div>
    ),
  },
  {
    key: 'replication',
    label: 'Replication',
    render: (database) => (
      <StatusBadge
        label={database.logicalReplication ? 'Enabled' : 'Disabled'}
        tone={database.logicalReplication ? 'info' : 'neutral'}
      />
    ),
  },
  {
    key: 'actions',
    label: '',
    render: (database) => (
      <Button asChild size="sm" variant="ghost">
        <Link to="/databases/$dbId" params={{ dbId: database.id }}>
          <ExternalLink className="size-4" />
          Open
        </Link>
      </Button>
    ),
  },
];

export const DatabaseList = () => {
  const databases = useListDatabases();

  if (databases.isLoading) {
    return (
      <div className="grid gap-2 rounded-lg border p-4">
        {Array.from({ length: 5 }).map((_, index) => (
          <Skeleton key={index.toString()} className="h-12 w-full" />
        ))}
      </div>
    );
  }

  if (databases.isError) {
    return (
      <EmptyPanel
        icon={ServerCrash}
        title="Database inventory is unavailable"
        description="The daemon returned an error while loading databases. Check migrations and keep the backend process running."
        actionLabel="Retry"
      />
    );
  }

  const rows = databases.data?.data ?? [];

  if (rows.length === 0) {
    return (
      <EmptyPanel
        icon={Database}
        title="No databases yet"
        description="Create a managed Postgres, MySQL, Redis, or compatible service for persistent state in your projects."
        actionLabel="New database"
      />
    );
  }

  return (
    <ResourceTable columns={databaseColumns} rows={rows} getRowKey={(database) => database.id} />
  );
};
