import type * as React from 'react';
import { cn } from '@/lib/utils';
import { SIDEBAR_WIDTH, SIDEBAR_WIDTH_ICON } from './sidebar-constants';

export function SidebarWrapper({ className, style, ...props }: React.ComponentProps<'div'>) {
  return (
    <div
      data-slot="sidebar-wrapper"
      style={
        {
          '--sidebar-width': SIDEBAR_WIDTH,
          '--sidebar-width-icon': SIDEBAR_WIDTH_ICON,
          ...style,
        } as React.CSSProperties
      }
      className={cn(
        'group/sidebar-wrapper flex min-h-svh w-full has-data-[variant=inset]:bg-sidebar',
        className
      )}
      {...props}
    />
  );
}
