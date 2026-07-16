import {
  Bell,
  Cloud,
  Code,
  Database,
  Download,
  FolderKanban,
  HardDrive,
  LayoutDashboard,
  LayoutTemplate,
  ScrollText,
  Settings,
  Terminal,
  Users,
} from 'lucide-react';
import { NavItem, type NavItemProps } from './nav-item';
import { UserMenu } from './user-menu';

type NavGroup = {
  title?: string;
  items: (NavItemProps & { exact?: boolean })[];
};

const navGroups: NavGroup[] = [
  {
    title: 'Overview',
    items: [
      { title: 'Dashboard', url: '/', icon: LayoutDashboard, exact: true },
      { title: 'Projects', url: '/projects', icon: FolderKanban },
    ],
  },
  {
    title: 'Resources',
    items: [
      { title: 'Databases', url: '/databases', icon: Database },
      { title: 'Storage', url: '/storage', icon: HardDrive },
      { title: 'Sources', url: '/settings/git-apps', icon: Code },
    ],
  },
  {
    title: 'Discover',
    items: [
      {
        title: 'Templates',
        url: '/templates',
        icon: LayoutTemplate,
      },
      { title: 'Import', url: '/imports/railway', icon: Download },
    ],
  },
  {
    title: 'System',
    items: [
      { title: 'Notifications', url: '/notifications', icon: Bell },
      { title: 'Audit Logs', url: '/audit-logs', icon: ScrollText },
      { title: 'Terminal', url: '/terminal', icon: Terminal },
<<<<<<< HEAD
      { title: 'Users', url: '/settings/users', icon: Users },
      { title: 'Settings', url: '/settings', icon: Settings },
=======
      { title: 'Settings', url: '/settings', icon: Settings, exact: true },
>>>>>>> ebe5d02 (feat: ui/ux revamp)
    ],
  },
];

const bottomNav = [
  {
    title: 'Docs',
    url: 'https://docs.vessl.com',
    icon: ScrollText,
    external: true,
  },
];

export function AppSidebar() {
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

      <nav className="flex flex-1 flex-col gap-6 overflow-y-auto px-3 py-4">
        {navGroups.map((group, i) => (
          <div key={i} className="flex flex-col gap-1">
            {group.title && (
              <h4 className="mb-1 px-2 font-medium text-[11px] text-sidebar-foreground/50 uppercase tracking-wider">
                {group.title}
              </h4>
            )}
            {group.items.map((item) => (
              <NavItem key={item.url} item={item} exact={item.exact} />
            ))}
          </div>
        ))}
      </nav>

      <div className="mt-auto flex flex-col gap-1 px-3 pt-4 pb-2">
        {bottomNav.map((item) => (
          <NavItem key={item.url} item={item} />
        ))}
      </div>

      <UserMenu />
    </aside>
  );
}
