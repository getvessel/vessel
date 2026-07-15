import {
  Bell,
  Cloud,
  Code,
  Database,
  FileText,
  HardDrive,
  Heart,
  LayoutDashboard,
  MessageSquare,
  Settings,
} from 'lucide-react';
import { NavItem } from './nav-item';
import { UserMenu } from './user-menu';

const mainNav = [
  { title: 'Dashboard', url: '/', icon: LayoutDashboard, exact: true },
  { title: 'Marketplace', url: '/marketplace', icon: Cloud },
  { title: 'Databases', url: '/databases', icon: Database },
  { title: 'S3 Storages', url: '/storages', icon: HardDrive },
  { title: 'Sources', url: '/sources', icon: Code },
  { title: 'Notifications', url: '/notifications', icon: Bell },
  { title: 'Audit Logs', url: '/audit-logs', icon: FileText },
];

const bottomNav = [
  {
    title: 'Docs',
    url: 'https://docs.vessl.com',
    icon: FileText,
    external: true,
  },
  {
    title: 'Feedback',
    url: 'https://feedback.vessl.com',
    icon: MessageSquare,
    external: true,
  },
  {
    title: 'Sponsor us',
    url: 'https://github.com/sponsors/vessl',
    icon: Heart,
    external: true,
  },
  { title: 'Settings', url: '/settings', icon: Settings },
];

export function AppSidebar() {
  const visibleMainNav = mainNav;
  const visibleBottomNav = bottomNav;

  return (
    <aside className="fixed inset-y-0 left-0 z-20 flex w-60 flex-col border-sidebar-border border-r bg-sidebar">
      <div className="flex h-14 shrink-0 items-center gap-2.5 border-sidebar-border border-b px-4">
        <div className="flex h-7 w-7 items-center justify-center rounded-lg bg-primary/10">
          <Cloud className="h-4 w-4 text-primary" />
        </div>
        <span className="font-semibold text-[15px] text-sidebar-foreground tracking-tight">
          Vessl
        </span>
        <span className="ml-auto rounded bg-sidebar-accent px-1.5 py-0.5 font-medium text-[10px] text-muted-foreground">
          v0.1
        </span>
      </div>

      <nav className="flex flex-1 flex-col gap-1 overflow-y-auto px-3 py-2">
        {visibleMainNav.map((item) => (
          <NavItem key={item.url} item={item} exact={item.exact} />
        ))}
      </nav>

      <div className="mt-auto flex flex-col gap-1 px-3 pb-2">
        {visibleBottomNav.map((item) => (
          <NavItem key={item.url} item={item} />
        ))}
      </div>

      <UserMenu />
    </aside>
  );
}
