import { createFileRoute, Outlet } from '@tanstack/react-router';
import { Shell } from '#/components/layout/shell';

export const Route = createFileRoute('/_dashboard')({
  component: DashboardLayout,
});

function DashboardLayout() {
  return (
    <Shell>
      <Outlet />
    </Shell>
  );
}
