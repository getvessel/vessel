import { useRouterState } from '@tanstack/react-router';
import type * as React from 'react';
import { useState } from 'react';
import { AppSidebar } from './app-sidebar';
import { CommandPalette } from './command-palette';
import { Topbar } from './topbar';

export function Shell({ children }: { children: React.ReactNode }) {
  const [commandOpen, setCommandOpen] = useState(false);
  const pathname = useRouterState({ select: (state) => state.location.pathname });

  return (
    <div className="flex min-h-screen bg-background">
      <AppSidebar />
<<<<<<< HEAD
      <div className="flex flex-1 flex-col pl-55">
        <Topbar />
=======
      <div className="flex flex-1 flex-col pl-60">
        <Topbar onOpenCommand={() => setCommandOpen(true)} />
>>>>>>> ebe5d02 (feat: ui/ux revamp)
        <main className="flex-1 overflow-auto p-6">
          <div key={pathname} className="page-transition mx-auto w-full max-w-[1480px]">
            {children}
          </div>
        </main>
      </div>
      <CommandPalette open={commandOpen} onOpenChange={setCommandOpen} />
    </div>
  );
}
