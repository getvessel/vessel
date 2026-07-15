import type { QueryClient } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { createRootRouteWithContext, HeadContent, Scripts } from '@tanstack/react-router';
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools';
import { ThemeProvider } from '#/components/theme-provider';
import { Toaster } from '#/components/ui/sonner';
import { TooltipProvider } from '#/components/ui/tooltip';
import appCss from '#/styles.css?url';

interface MyRouterContext {
  queryClient: QueryClient;
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  head: () => ({
    meta: [
      {
        charSet: 'utf-8',
      },
      {
        name: 'viewport',
        content: 'width=device-width, initial-scale=1',
      },
      {
        title: 'Vessl Dashboard',
      },
    ],
    links: [
      {
        rel: 'stylesheet',
        href: appCss,
      },
    ],
  }),
  shellComponent: RootDocument,
});

import { useEffect } from 'react';
import { useGetPublicSettings } from '#/hooks/useSettings';

function DynamicTitle() {
  const { data } = useGetPublicSettings();
  const siteName = data?.data?.siteName;

  useEffect(() => {
    if (siteName) {
      document.title = `${siteName} Dashboard`;
    } else {
      document.title = 'Vessl Dashboard';
    }
  }, [siteName]);

  return null;
}

function RootDocument({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        <HeadContent />
      </head>
      <body>
        <ThemeProvider attribute="class" defaultTheme="dark" enableSystem disableTransitionOnChange>
          <DynamicTitle />
          <TooltipProvider>{children}</TooltipProvider>
          <Toaster />
          <TanStackRouterDevtools position="bottom-right" />
          <ReactQueryDevtools position="bottom-left" />
        </ThemeProvider>
        <Scripts />
      </body>
    </html>
  );
}
