import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_dashboard/terminal')({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Terminal</div>;
}
