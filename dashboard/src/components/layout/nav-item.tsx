import { Link, useRouterState } from '@tanstack/react-router';
import { ExternalLink } from 'lucide-react';
import type React from 'react';
import { useInterfaceSound } from '../../hooks/use-interface-sound';

export type NavItemProps = {
  title: string;
  url: string;
  icon: React.ComponentType<{ className?: string }>;
  external?: boolean;
  badge?: string;
};

export function NavItem({ item, exact = false }: { item: NavItemProps; exact?: boolean }) {
  const routerState = useRouterState();
  const { playNavigationSound } = useInterfaceSound();
  const pathname = routerState.location.pathname;
  const isActive = exact
    ? pathname === item.url
    : pathname.startsWith(item.url) && item.url !== '/';

  const className = [
    'group relative flex min-h-9 items-center gap-2 rounded-md px-2 py-1.5 text-[13px] font-medium outline-none transition-[background-color,color,box-shadow,transform] duration-150 ease-out focus-visible:ring-2 focus-visible:ring-sidebar-ring/50 active:scale-[0.985]',
    isActive
      ? 'bg-sidebar-accent text-sidebar-accent-foreground shadow-[inset_0_0_0_1px_rgb(255_255_255/0.04)]'
      : 'text-sidebar-foreground/60 hover:bg-sidebar-accent/70 hover:text-sidebar-foreground',
  ].join(' ');

  const IconComponent = (
    <span
      className={[
        'flex h-6 w-6 shrink-0 items-center justify-center rounded-md border transition-[background-color,border-color,color,transform] duration-150',
        isActive
          ? 'border-primary/25 bg-primary/10 text-primary'
          : 'border-sidebar-border/70 bg-sidebar-accent/25 text-sidebar-foreground/50 group-hover:scale-[1.03] group-hover:border-sidebar-border group-hover:bg-sidebar-accent group-hover:text-sidebar-foreground/80',
      ].join(' ')}
    >
      <item.icon className="h-3.5 w-3.5" />
    </span>
  );

  const ActiveIndicator = isActive ? (
    <span className="absolute top-1/2 -left-3 h-4 w-0.5 -translate-y-1/2 rounded-full bg-primary" />
  ) : null;

  if (item.external) {
    return (
      <a
        href={item.url}
        target="_blank"
        rel="noopener noreferrer"
        className={className}
        onClick={playNavigationSound}
      >
        {ActiveIndicator}
        {IconComponent}
        <span className="flex-1 truncate">{item.title}</span>
        {item.badge && (
          <span className="ml-auto flex h-5 items-center justify-center rounded-full bg-primary/10 px-1.5 font-medium text-[10px] text-primary">
            {item.badge}
          </span>
        )}
        <ExternalLink className="h-3.5 w-3.5 shrink-0 opacity-50" />
      </a>
    );
  }

  return (
    <Link to={item.url as never} className={className} onClick={playNavigationSound}>
      {ActiveIndicator}
      {IconComponent}
      <span className="flex-1 truncate">{item.title}</span>
      {item.badge && (
        <span className="ml-auto flex h-5 items-center justify-center rounded-full bg-primary/10 px-1.5 font-medium text-[10px] text-primary">
          {item.badge}
        </span>
      )}
    </Link>
  );
}
