import { createFileRoute, Outlet, redirect } from '@tanstack/react-router';
import { Shell } from '#/components/layout/shell';
import { authStore } from '#/stores/authStore';

export const Route = createFileRoute('/_shell')({
  beforeLoad: () => {
    if (!authStore.state.isAuthenticated) {
      throw redirect({ to: '/login' });
    }
  },
  component: ShellLayout,
});

function ShellLayout() {
  return (
    <Shell>
      <Outlet />
    </Shell>
  );
}
