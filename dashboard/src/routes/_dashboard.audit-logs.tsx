import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/audit-logs')({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Audit Logs</div>;
}
