import {
  createFileRoute,
  Outlet,
  redirect,
  useNavigate,
  useRouterState,
} from '@tanstack/react-router';
import { useStore } from '@tanstack/react-store';
import { useEffect } from 'react';
import { Shell } from '#/components/layout/shell';
import { useGetPublicSettings } from '#/hooks/useSettings';
import { useListWorkspaces } from '#/hooks/useWorkspaces';
import { authStore } from '#/stores/authStore';
import { workspaceStore } from '#/stores/workspaceStore';

export const Route = createFileRoute('/_workspace')({
  beforeLoad: () => {
    if (!authStore.state.isAuthenticated) {
      throw redirect({ to: '/login' });
    }
  },
  component: WorkspaceLayout,
});

function WorkspaceLayout() {
  useListWorkspaces();
  const router = useRouterState();
  const navigate = useNavigate();
  const { data: settings } = useGetPublicSettings();
  const { activeWorkspace } = useStore(workspaceStore);

  const isCloudMode = settings?.data?.isCloudMode;
  const isSubscribed = activeWorkspace?.subscriptionStatus === 'active';
  const requiresSubscription = isCloudMode && !isSubscribed && activeWorkspace;
  const isSubscriptionRoute = router.location.pathname.startsWith('/subscribe');

  useEffect(() => {
    if (requiresSubscription && !isSubscriptionRoute) {
      navigate({ to: '/subscribe' });
    }
  }, [requiresSubscription, isSubscriptionRoute, navigate]);

  return (
    <Shell>
      <Outlet />
    </Shell>
  );
}
