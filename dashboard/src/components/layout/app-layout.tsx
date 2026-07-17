import { useRouterState } from '@tanstack/react-router';
import type * as React from 'react';
import { useEffect, useState } from 'react';
import { AppSidebar } from './app-sidebar';
import { BackgroundPattern } from './background-pattern';
import { CommandPalette } from './command-palette';
import { Topbar } from './topbar';

export function AppLayout({ children }: { children: React.ReactNode }) {
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [commandOpen, setCommandOpen] = useState(false);
  const pathname = useRouterState({
    select: (state) => state.location.pathname,
  });

  // biome-ignore lint/correctness/useExhaustiveDependencies: run on pathname change
  useEffect(() => {
    setMobileMenuOpen(false);
  }, [pathname]);

  return (
    <div className="relative flex min-h-screen bg-background">
      <BackgroundPattern />
      <AppSidebar
        collapsed={sidebarCollapsed}
        onToggle={() => setSidebarCollapsed((p) => !p)}
        mobileOpen={mobileMenuOpen}
        onMobileClose={() => setMobileMenuOpen(false)}
      />
      <div
        className={`relative flex flex-1 flex-col ${sidebarCollapsed ? 'md:pl-16' : 'md:pl-64'}`}
      >
        <Topbar
          onOpenCommand={() => setCommandOpen(true)}
          onMenuToggle={() => setMobileMenuOpen((p) => !p)}
          mobileMenuOpen={mobileMenuOpen}
        />
        <main className="flex-1 overflow-auto p-4 md:p-8">
          <div key={pathname} className="page-transition mx-auto w-full max-w-7xl">
            {children}
          </div>
        </main>
      </div>
      <CommandPalette open={commandOpen} onOpenChange={setCommandOpen} />
    </div>
  );
}
