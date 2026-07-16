import { createFileRoute, Outlet } from '@tanstack/react-router';
import { AuthLayout } from '#/features/auth/auth-layout';

export const Route = createFileRoute('/_auth')({
  component: AuthRoute,
});

function AuthRoute() {
  return (
    <AuthLayout>
      <Outlet />
    </AuthLayout>
  );
}
