import { createFileRoute, Outlet, redirect } from '@tanstack/react-router';
import { Layers } from 'lucide-react';

import { authStore } from '#/stores/authStore';

export const Route = createFileRoute('/_auth')({
  beforeLoad: () => {
    if (authStore.state.isAuthenticated) {
      throw redirect({ to: '/' });
    }
  },
  component: AuthLayout,
});

function AuthLayout() {
  return (
    <div className="flex min-h-screen w-full flex-col items-center justify-center bg-background p-4 font-sans text-foreground">
      <div className="relative z-10 flex w-full max-w-lg flex-col items-center sm:max-w-120">
        <div className="mb-8 flex items-center justify-center gap-2 font-bold text-xl">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
            <Layers className="h-5 w-5" />
          </div>
          Vessl
        </div>

        <div className="w-full">
          <Outlet />
        </div>
      </div>
    </div>
  );
}
