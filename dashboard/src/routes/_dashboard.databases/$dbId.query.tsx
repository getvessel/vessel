import { createFileRoute } from '@tanstack/react-router';
import { Play, Rows3 } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { Textarea } from '#/components/ui/textarea';
import { OperationalPage } from '#/features/dashboard/operational-page';
import { ResourceTable } from '#/features/dashboard/resource-table';
import { useGetSchemas } from '#/hooks/useDatabases';

export const Route = createFileRoute('/_dashboard/databases/$dbId/query')({
  component: DatabaseQueryPage,
});

function DatabaseQueryPage() {
  const { dbId } = Route.useParams();
  const schemas = useGetSchemas(dbId);
  const rows = schemas.data?.data ?? [];

  return (
    <OperationalPage
      title="SQL query"
      description="Run inspected SQL against this database and review schemas before changing data."
      scope={dbId}
      statusLabel={schemas.isError ? 'Schema unavailable' : 'Ready'}
      statusTone={schemas.isError ? 'danger' : 'healthy'}
      primaryAction={{ label: 'Run query', icon: Play }}
      metrics={[
        {
          label: 'Tables',
          value: schemas.isLoading ? undefined : rows.length.toString(),
          detail: 'Visible schemas',
          icon: Rows3,
        },
      ]}
    >
      <div className="grid gap-4 xl:grid-cols-[1fr_360px]">
        <div className="rounded-lg border bg-card">
          <div className="border-b p-3">
            <Textarea
              defaultValue="select * from information_schema.tables limit 20;"
              className="min-h-52 resize-none border-0 bg-[linear-gradient(to_right,var(--border)_1px,transparent_1px),linear-gradient(to_bottom,var(--border)_1px,transparent_1px)] bg-[size:32px_32px] font-mono text-sm shadow-none focus-visible:ring-0"
            />
          </div>
          <div className="flex justify-end p-3">
            <Button>
              <Play className="size-4" />
              Execute
            </Button>
          </div>
        </div>
        <ResourceTable
          columns={[
            { key: 'name', label: 'Table', render: (table) => table.name },
            { key: 'columns', label: 'Columns', render: (table) => table.columns.length },
          ]}
          rows={rows}
          getRowKey={(table) => table.name}
        />
      </div>
    </OperationalPage>
  );
}
