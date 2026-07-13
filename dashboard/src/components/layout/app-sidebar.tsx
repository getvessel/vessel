import {
  ArrowUpRight01Icon,
  AuditIcon,
  BellIcon,
  BucketIcon,
  ChartHistogramIcon,
  CheckmarkCircle01Icon,
  ChevronDownIcon,
  ClipboardListIcon,
  CloudServerIcon,
  DashboardSquare01Icon,
  Database01Icon,
  FavouriteIcon,
  HelpCircleIcon,
  Logout01Icon,
  Message01Icon,
  Moon02Icon,
  MoreVerticalIcon,
  PlusSignIcon,
  Settings01Icon,
  SourceCodeIcon,
  Sun01Icon,
  UserCircleIcon,
  UserMultiple02Icon,
} from '@hugeicons/core-free-icons';
import { HugeiconsIcon } from '@hugeicons/react';
import { Link, useRouterState } from '@tanstack/react-router';
import { useStore } from '@tanstack/react-store';
import { useTheme } from 'next-themes';
import { useEffect, useMemo, useRef, useState } from 'react';
import { useLogout } from '#/hooks/useAuth';
import { useListWorkspaces } from '#/hooks/useWorkspaces';
import type { Workspace } from '#/interfaces/workspace';
import { authStore } from '#/stores/authStore';
import { workspaceActions, workspaceStore } from '#/stores/workspaceStore';

// ─── Nav structure ────────────────────────────────────────────────────────────

const mainNav = [
  { title: 'Dashboard', url: '/', icon: DashboardSquare01Icon, exact: true },
  { title: 'Workspaces', url: '/workspaces', icon: UserMultiple02Icon },
  { title: 'Databases', url: '/databases', icon: Database01Icon },
  { title: 'S3 Storages', url: '/storages', icon: BucketIcon },
  { title: 'Sources', url: '/sources', icon: SourceCodeIcon },
  { title: 'Notifications', url: '/notifications', icon: BellIcon },
  { title: 'Audit Logs', url: '/audit-logs', icon: AuditIcon },
];

const bottomNav = [
  {
    title: 'Docs',
    url: 'https://docs.vessl.com',
    icon: ClipboardListIcon,
    external: true,
  },
  {
    title: 'Feedback',
    url: 'https://feedback.vessl.com',
    icon: Message01Icon,
    external: true,
  },
  {
    title: 'Sponsor us',
    url: 'https://github.com/sponsors/vessl',
    icon: FavouriteIcon,
    external: true,
  },
  { title: 'Settings', url: '/settings', icon: Settings01Icon },
];

// ─── Nav item ─────────────────────────────────────────────────────────────────

function NavItem({
  item,
  exact = false,
}: {
  item: { title: string; url: string; icon: unknown; external?: boolean };
  exact?: boolean;
}) {
  const routerState = useRouterState();
  const pathname = routerState.location.pathname;
  const isActive = exact
    ? pathname === item.url
    : pathname.startsWith(item.url) && item.url !== '/';

  const className = [
    'group flex items-center gap-3 rounded-md px-3 py-2 text-[13px] font-medium transition-all duration-100',
    isActive
      ? 'bg-sidebar-accent text-sidebar-accent-foreground'
      : 'text-sidebar-foreground/60 hover:text-sidebar-foreground hover:bg-sidebar-accent/60',
  ].join(' ');

  const IconComponent = (
    <HugeiconsIcon
      // @ts-expect-error icon type is opaque
      icon={item.icon}
      className={[
        'h-[16px] w-[16px] shrink-0',
        isActive ? 'text-primary' : 'group-hover:text-sidebar-foreground/70',
      ].join(' ')}
    />
  );

  if (item.external) {
    return (
      <a href={item.url} target="_blank" rel="noopener noreferrer" className={className}>
        {IconComponent}
        <span className="truncate flex-1">{item.title}</span>
        <HugeiconsIcon icon={ArrowUpRight01Icon} className="h-3.5 w-3.5 opacity-50 shrink-0" />
      </a>
    );
  }

  return (
    <Link to={item.url as never} className={className}>
      {IconComponent}
      <span className="truncate">{item.title}</span>
    </Link>
  );
}

// ─── Workspace switcher ────────────────────────────────────────────────────────

function WorkspaceSwitcher() {
  const { data } = useListWorkspaces();
  const workspaceState = useStore(workspaceStore);
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  const workspaces: Workspace[] = useMemo(
    () => (data as { data?: Workspace[] } | undefined)?.data ?? [],
    [data]
  );
  const active = workspaceState.activeWorkspace;

  // Sync fetched list into store
  useEffect(() => {
    if (workspaces.length > 0) {
      workspaceActions.setWorkspaces(workspaces);
    }
  }, [workspaces]);

  // Close on outside click
  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener('mousedown', handler);
    return () => document.removeEventListener('mousedown', handler);
  }, []);

  const initials = active?.name ? active.name.slice(0, 2).toUpperCase() : '??';

  return (
    <div ref={ref} className="relative px-3 pb-2">
      <button
        type="button"
        onClick={() => setOpen((o) => !o)}
        className="flex w-full items-center gap-2.5 rounded-lg border border-sidebar-border bg-sidebar-accent/40 px-2.5 py-2 text-left hover:bg-sidebar-accent/70 transition-colors duration-100"
      >
        {/* Workspace avatar */}
        <div className="flex h-6 w-6 shrink-0 items-center justify-center rounded-md bg-primary/20 text-[10px] font-bold text-primary">
          {initials}
        </div>
        <div className="min-w-0 flex-1">
          <p className="truncate text-[12px] font-semibold text-sidebar-foreground leading-none">
            {active?.name ?? 'Select workspace'}
          </p>
          {active?.preferredRegion && (
            <p className="truncate text-[10px] text-sidebar-foreground/40 mt-0.5 leading-none">
              {active.preferredRegion}
            </p>
          )}
        </div>
        <HugeiconsIcon
          icon={ChevronDownIcon}
          className={[
            'h-3.5 w-3.5 shrink-0 text-sidebar-foreground/40 transition-transform duration-150',
            open ? 'rotate-180' : '',
          ].join(' ')}
        />
      </button>

      {open && (
        <div className="absolute left-3 right-3 top-full z-50 mt-1 rounded-lg border border-border bg-popover shadow-lg overflow-hidden">
          <div className="p-1 max-h-48 overflow-y-auto">
            {workspaces.length === 0 && (
              <p className="px-2 py-3 text-center text-xs text-muted-foreground">
                No workspaces yet
              </p>
            )}
            {workspaces.map((ws) => (
              <button
                key={ws.id}
                type="button"
                onClick={() => {
                  workspaceActions.switchWorkspace(ws);
                  setOpen(false);
                }}
                className="flex w-full items-center gap-2.5 rounded-md px-2 py-1.5 text-sm hover:bg-accent transition-colors"
              >
                <div className="flex h-5 w-5 shrink-0 items-center justify-center rounded text-[9px] font-bold bg-primary/15 text-primary">
                  {ws.name.slice(0, 2).toUpperCase()}
                </div>
                <span className="flex-1 truncate text-[12px] font-medium text-foreground">
                  {ws.name}
                </span>
                {ws.id === active?.id && (
                  <HugeiconsIcon
                    icon={CheckmarkCircle01Icon}
                    className="h-3.5 w-3.5 text-primary shrink-0"
                  />
                )}
              </button>
            ))}
          </div>
          <div className="border-t border-border p-1">
            <Link
              to={'/workspaces' as never}
              onClick={() => setOpen(false)}
              className="flex w-full items-center gap-2 rounded-md px-2 py-1.5 text-xs font-medium text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
            >
              <HugeiconsIcon icon={PlusSignIcon} className="h-3.5 w-3.5" />
              New workspace
            </Link>
          </div>
        </div>
      )}
    </div>
  );
}

// ─── User Profile Menu ────────────────────────────────────────────────────────

function UserMenu() {
  const { theme, setTheme } = useTheme();
  const authState = useStore(authStore);
  const { mutateAsync: logout } = useLogout();
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  const user = authState.user;
  const initials = user?.name
    ? user.name
        .split(' ')
        .map((n) => n[0])
        .join('')
        .toUpperCase()
    : 'U';

  // Close on outside click
  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener('mousedown', handler);
    return () => document.removeEventListener('mousedown', handler);
  }, []);

  return (
    <div ref={ref} className="relative border-t border-sidebar-border p-3">
      {open && (
        <div className="absolute bottom-full left-3 right-3 mb-2 rounded-lg border border-border bg-popover shadow-xl overflow-hidden py-1 z-50 animate-in fade-in zoom-in-95 duration-100">
          <div className="px-3 py-2 border-b border-border mb-1">
            <p className="text-[13px] font-semibold text-foreground">{user?.name}</p>
            <p className="text-[11px] text-muted-foreground">{user?.email}</p>
          </div>

          <div className="px-1 py-1 space-y-0.5">
            <Link
              to={'/settings' as never}
              onClick={() => setOpen(false)}
              className="flex items-center gap-2.5 rounded-md px-2 py-1.5 text-[13px] hover:bg-accent transition-colors"
            >
              <HugeiconsIcon icon={UserCircleIcon} className="h-4 w-4 text-muted-foreground" />
              Account Settings
            </Link>
            <Link
              to={'/settings' as never}
              onClick={() => setOpen(false)}
              className="flex items-center gap-2.5 rounded-md px-2 py-1.5 text-[13px] hover:bg-accent transition-colors"
            >
              <HugeiconsIcon icon={Settings01Icon} className="h-4 w-4 text-muted-foreground" />
              Workspace Settings
            </Link>
            <Link
              to={'/settings' as never}
              onClick={() => setOpen(false)}
              className="flex items-center gap-2.5 rounded-md px-2 py-1.5 text-[13px] hover:bg-accent transition-colors"
            >
              <HugeiconsIcon icon={ChartHistogramIcon} className="h-4 w-4 text-muted-foreground" />
              Project Usage
            </Link>
          </div>

          <div className="my-1 h-px bg-border" />

          <div className="px-1 py-1 space-y-0.5">
            <Link
              to={'/docs' as never}
              onClick={() => setOpen(false)}
              className="flex items-center gap-2.5 rounded-md px-2 py-1.5 text-[13px] hover:bg-accent transition-colors"
            >
              <HugeiconsIcon icon={ClipboardListIcon} className="h-4 w-4 text-muted-foreground" />
              Documentation
            </Link>
            <Link
              to={'/support' as never}
              onClick={() => setOpen(false)}
              className="flex items-center gap-2.5 rounded-md px-2 py-1.5 text-[13px] hover:bg-accent transition-colors"
            >
              <HugeiconsIcon icon={HelpCircleIcon} className="h-4 w-4 text-muted-foreground" />
              Support
            </Link>
          </div>

          <div className="my-1 h-px bg-border" />

          <div className="px-1 py-1 space-y-0.5">
            <button
              type="button"
              onClick={() => {
                setTheme(theme === 'dark' ? 'light' : 'dark');
                setOpen(false);
              }}
              className="flex w-full items-center gap-2.5 rounded-md px-2 py-1.5 text-[13px] hover:bg-accent transition-colors"
            >
              <HugeiconsIcon
                icon={theme === 'dark' ? Sun01Icon : Moon02Icon}
                className="h-4 w-4 text-muted-foreground"
              />
              {theme === 'dark' ? 'Light Theme' : 'Dark Theme'}
            </button>
            <button
              type="button"
              onClick={() => logout()}
              className="flex w-full items-center gap-2.5 rounded-md px-2 py-1.5 text-[13px] text-destructive hover:bg-destructive/10 transition-colors"
            >
              <HugeiconsIcon icon={Logout01Icon} className="h-4 w-4" />
              Log out
            </button>
          </div>
        </div>
      )}

      <button
        type="button"
        onClick={() => setOpen(!open)}
        className="flex w-full items-center gap-2.5 rounded-lg hover:bg-sidebar-accent/50 p-1.5 transition-colors"
      >
        <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-primary/20 text-[11px] font-bold text-primary">
          {initials}
        </div>
        <div className="min-w-0 text-left flex-1">
          <p className="truncate text-[12px] font-semibold text-sidebar-foreground leading-none">
            {user?.name ?? 'User'}
          </p>
          <p className="truncate text-[10px] text-sidebar-foreground/50 mt-0.5 leading-none">
            Workspace Owner
          </p>
        </div>
        <HugeiconsIcon
          icon={MoreVerticalIcon}
          className="h-4 w-4 text-sidebar-foreground/50 shrink-0"
        />
      </button>
    </div>
  );
}

// ─── Sidebar ──────────────────────────────────────────────────────────────────

export function AppSidebar() {
  return (
    <aside className="fixed inset-y-0 left-0 z-20 flex w-[240px] flex-col border-r border-sidebar-border bg-sidebar">
      {/* Logo */}
      <div className="flex h-14 shrink-0 items-center gap-2.5 px-4 border-b border-sidebar-border">
        <div className="flex h-7 w-7 items-center justify-center rounded-lg bg-primary/10">
          <HugeiconsIcon icon={CloudServerIcon} className="h-4 w-4 text-primary" />
        </div>
        <span className="font-semibold text-[15px] tracking-tight text-sidebar-foreground">
          Vessl
        </span>
        <span className="ml-auto text-[10px] font-medium text-muted-foreground bg-sidebar-accent px-1.5 py-0.5 rounded">
          v0.1
        </span>
      </div>

      {/* Workspace switcher */}
      <div className="shrink-0 pt-4">
        <WorkspaceSwitcher />
      </div>

      {/* Main nav */}
      <nav className="flex flex-1 flex-col overflow-y-auto px-3 py-2 gap-1">
        {mainNav.map((item) => (
          <NavItem key={item.url} item={item} exact={item.exact} />
        ))}
      </nav>

      {/* Bottom Nav */}
      <div className="mt-auto flex flex-col px-3 pb-2 gap-1">
        {bottomNav.map((item) => (
          <NavItem key={item.url} item={item} />
        ))}
      </div>

      <UserMenu />
    </aside>
  );
}
