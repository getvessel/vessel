import { createFileRoute } from '@tanstack/react-router';
import { Rows3, Table2 } from 'lucide-react';
import { EmptyPanel } from '#/features/dashboard/empty-panel';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { ResourceTable } from '#/features/dashboard/resource-table';
import { useGetSchemas } from '#/hooks/useDatabases';

export const Route = createFileRoute('/_dashboard/databases/$dbId/data')({
  component: DatabaseDataPage,
});

function DatabaseDataPage() {
  const { dbId } = Route.useParams();
  const schemas = useGetSchemas(dbId);
  const rows = schemas.data?.data ?? [];

  return (
    <OperationalPage
      title="Data browser"
      description="Browse tables, inspect columns, and open row data without leaving the database context."
      scope={dbId}
      statusLabel={schemas.isError ? 'Schema error' : 'Browsing'}
      statusTone={schemas.isError ? 'danger' : 'healthy'}
      metrics={[
        {
          label: 'Tables',
          value: schemas.isLoading ? undefined : rows.length.toString(),
          detail: 'Detected in schema',
          icon: Table2,
        },
      ]}
    >
      {rows.length > 0 ? (
        <ResourceTable
          columns={[
            { key: 'name', label: 'Table', render: (table) => table.name },
            { key: 'columns', label: 'Columns', render: (table) => table.columns.length },
            {
              key: 'primary',
              label: 'Primary keys',
              render: (table) => table.columns.filter((column) => column.isPrimary).length,
            },
          ]}
          rows={rows}
          getRowKey={(table) => table.name}
        />
      ) : (
        <EmptyPanel
          icon={Rows3}
          title="No tables detected"
          description="Create tables or import data, then return here to browse rows and inspect schema structure."
          actionLabel="Open SQL query"
        />
      )}
    </OperationalPage>
  );
}
