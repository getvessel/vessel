import { Link, Outlet, useRouterState } from '@tanstack/react-router';
import { Bell, Lock, Settings } from 'lucide-react';

type Tab = { id: string; label: string; icon: React.ReactNode; path: string };

const TABS: Tab[] = [
  {
    id: 'general',
    label: 'General',
    icon: <Settings className="h-4 w-4" />,
    path: '/settings/general',
  },
  {
    id: 'notifications',
    label: 'Notifications',
    icon: <Bell className="h-4 w-4" />,
    path: '/settings/notifications',
  },
  { id: 'oauth', label: 'OAuth', icon: <Lock className="h-4 w-4" />, path: '/settings/oauth' },
];

export const SettingsLayout = () => {
  const { location } = useRouterState();
  const pathname = location.pathname;

  const activeId = TABS.find((t) => pathname.startsWith(t.path))?.id ?? 'general';

  return (
    <div className="flex min-h-full flex-col">
      <div className="border-border/50 border-b bg-background/50 px-6 pt-6 backdrop-blur-sm">
        <div className="mb-5 flex items-center gap-3">
          <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-primary/10 text-primary">
            <Settings className="h-4.5 w-4.5" />
          </div>
          <div>
            <h1 className="font-bold text-xl">Instance Settings</h1>
            <p className="text-muted-foreground text-sm">Manage your Vessl server configuration</p>
          </div>
        </div>

        <nav className="flex gap-1 overflow-x-auto">
          {TABS.map((t) => {
            const isActive = t.id === activeId;
            return (
              <Link
                key={t.id}
                to={t.path}
                className={[
                  'flex shrink-0 items-center gap-2 rounded-t-lg border border-b-0 px-4 py-2.5 text-sm transition-colors',
                  isActive
                    ? 'border-border/50 bg-card font-medium text-foreground'
                    : 'border-transparent text-muted-foreground hover:text-foreground',
                ].join(' ')}
              >
                {t.icon}
                {t.label}
              </Link>
            );
          })}
        </nav>
      </div>

      <div className="flex-1 p-6">
        <Outlet />
      </div>
    </div>
  );
};
