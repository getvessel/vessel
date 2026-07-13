import { createFileRoute, Outlet, redirect } from '@tanstack/react-router';
import { Shell } from '#/components/layout/shell';
import { authStore } from '#/stores/authStore';

export const Route = createFileRoute('/_project')({
  beforeLoad: () => {
    if (!authStore.state.isAuthenticated) {
      throw redirect({ to: '/login' });
    }
  },
  component: () => (
    <Shell>
      <Outlet />
    </Shell>
  ),
});
