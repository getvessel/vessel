import { useNavigate, useRouterState } from '@tanstack/react-router';
import { BellIcon, PlusIcon, SearchIcon } from 'lucide-react';
import { toast } from 'sonner';

const routeLabels: Record<string, { title: string; description: string }> = {
  '/': {
    title: 'Instance overview',
    description: 'Health and resource pressure',
  },
  '/projects': { title: 'Projects', description: 'Workloads and environments' },
  '/databases': {
    title: 'Databases',
    description: 'Persistent services and backups',
  },
  '/templates': {
    title: 'Templates',
    description: 'Starter services and resources',
  },
  '/teams': { title: 'Teams', description: 'Ownership and permissions' },
  '/notifications': {
    title: 'Notifications',
    description: 'Event delivery channels',
  },
  '/audit-logs': {
    title: 'Audit logs',
    description: 'Security and operator activity',
  },
  '/terminal': { title: 'Terminal', description: 'Controlled host sessions' },
  '/settings': { title: 'Settings', description: 'Instance administration' },
};

const sectionDescriptions: Record<string, string> = {
  projects: 'Project workspace',
  services: 'Service operations',
  databases: 'Database operations',
  settings: 'Instance administration',
  imports: 'Migration workflow',
  profile: 'Account settings',
};

function titleFromPath(pathname: string) {
  const exact = routeLabels[pathname];

  if (exact) {
    return exact;
  }

  const segments = pathname.split('/').filter(Boolean);
  const section = segments[0] ?? 'dashboard';
  const current = segments.at(-1) ?? section;
  const title = current.replaceAll('-', ' ').replace(/^\w/, (character) => character.toUpperCase());

  return {
    title,
    description: sectionDescriptions[section] ?? 'Dashboard workflow',
  };
}

interface TopbarProps {
  onOpenCommand: () => void;
}

export function Topbar({ onOpenCommand }: TopbarProps) {
  const routerState = useRouterState();
  const navigate = useNavigate();
  const pathname = routerState.location.pathname;
  const current = titleFromPath(pathname);

  return (
    <header className="flex h-16 shrink-0 items-center justify-between bg-background/80 px-8 shadow-[inset_0_-1px_0_0_rgb(255_255_255_/_0.06)] backdrop-blur-xl">
      <div className="min-w-0">
        <h1 className="truncate font-semibold text-foreground text-lg tracking-tight">
          {current.title}
        </h1>
        {current.description && (
          <p className="text-muted-foreground text-xs">{current.description}</p>
        )}
      </div>

      <div className="flex items-center gap-2">
        <button
          type="button"
          onClick={onOpenCommand}
          className="flex h-9 items-center gap-2 rounded-xl border border-border/60 bg-muted/40 px-3 text-muted-foreground text-sm transition-all hover:border-border hover:bg-muted hover:text-foreground active:scale-[0.97]"
        >
          <SearchIcon className="h-4 w-4 shrink-0" />
          <span className="hidden sm:inline">Search...</span>
          <kbd className="rounded-md border bg-background/60 px-1.5 py-0.5 font-mono text-[11px] leading-none">
            ⌘K
          </kbd>
        </button>

        <button
          type="button"
          onClick={() =>
            toast.info('New resource', {
              description: 'Creation menu coming soon',
            })
          }
          className="flex h-9 items-center gap-1.5 rounded-xl bg-primary px-4 font-semibold text-primary-foreground text-sm shadow-lg shadow-primary/25 transition-all hover:brightness-110 active:scale-[0.97]"
        >
          <PlusIcon className="h-4 w-4" />
          <span>New</span>
        </button>

        <button
          type="button"
          onClick={() => navigate({ to: '/notifications' })}
          className="relative flex h-9 w-9 items-center justify-center rounded-xl border border-border/60 transition-colors hover:bg-muted"
        >
          <BellIcon className="h-4 w-4" />
          <div className="absolute top-2 right-2 h-1.5 w-1.5 rounded-full bg-primary ring-2 ring-background" />
        </button>
      </div>
    </header>
  );
}
