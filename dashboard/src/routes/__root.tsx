import type { QueryClient } from '@tanstack/react-query';
import { createRootRouteWithContext, Outlet } from '@tanstack/react-router';
import { DynamicTitle } from '#/components/layout/dynamic-title';
import { ThemeProvider } from '#/components/theme-provider';
import { Toaster } from '#/components/ui/sonner';
import { TooltipProvider } from '#/components/ui/tooltip';

interface MyRouterContext {
  queryClient: QueryClient;
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  component: RootDocument,
});

function RootDocument() {
  return (
    <ThemeProvider>
      <DynamicTitle />
      <TooltipProvider>
        <Outlet />
      </TooltipProvider>
      <Toaster />
    </ThemeProvider>
  );
}
