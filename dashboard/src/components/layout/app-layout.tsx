import { useRouterState } from '@tanstack/react-router';
import type * as React from 'react';
import { useState } from 'react';
import { AppSidebar } from './app-sidebar';
import { CommandPalette } from './command-palette';
import { Topbar } from './topbar';

export function AppLayout({ children }: { children: React.ReactNode }) {
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [commandOpen, setCommandOpen] = useState(false);
  const pathname = useRouterState({
    select: (state) => state.location.pathname,
  });

  return (
    <div className="flex min-h-screen bg-background">
      <AppSidebar collapsed={sidebarCollapsed} onToggle={() => setSidebarCollapsed((p) => !p)} />
      <div className={`flex flex-1 flex-col ${sidebarCollapsed ? 'pl-16' : 'pl-64'}`}>
        <Topbar onOpenCommand={() => setCommandOpen(true)} />
        <main className="flex-1 overflow-auto p-8">
          <div key={pathname} className="page-transition mx-auto w-full max-w-[1280px]">
            {children}
          </div>
        </main>
      </div>
      <CommandPalette open={commandOpen} onOpenChange={setCommandOpen} />
    </div>
  );
}
