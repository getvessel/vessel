import { Link } from '@tanstack/react-router';
import { BarChart3, LogOut, Moon, MoreVertical, Sun, UserCircle } from 'lucide-react';
import { useTheme } from 'next-themes';
import { useEffect, useRef, useState } from 'react';
import { useLogout } from '#/hooks/useAuth';
import { useAuthState } from '#/stores/authStore';

interface UserMenuProps {
  collapsed: boolean;
}

export function UserMenu({ collapsed }: UserMenuProps) {
  const { theme, setTheme } = useTheme();
  const authState = useAuthState();
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

  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener('mousedown', handler);
    return () => document.removeEventListener('mousedown', handler);
  }, []);

  return (
    <div
      ref={ref}
      className={`relative bg-sidebar-accent/20 ${collapsed ? 'px-1 py-1' : 'px-2 py-2'}`}
    >
      {open && (
        <div
          className={`fade-in zoom-in-95 absolute z-50 mb-2 animate-in rounded-xl border border-border/60 bg-popover/95 py-1 shadow-2xl backdrop-blur-2xl duration-100 ${
            collapsed ? 'bottom-full left-full ml-3 min-w-56' : 'right-2 bottom-full left-2'
          }`}
        >
          <div className="mb-1 border-border/50 border-b px-3 py-2">
            <p className="font-semibold text-[13px] text-foreground">{user?.name}</p>
            <p className="text-[11px] text-muted-foreground">{user?.email}</p>
          </div>

          <div className="space-y-0.5 px-1 py-0.5">
            <Link
              to={'/settings' as never}
              onClick={() => setOpen(false)}
              className="flex items-center gap-2.5 rounded-lg px-2 py-1.5 text-[13px] transition-colors hover:bg-accent/50"
            >
              <UserCircle className="h-4 w-4 text-muted-foreground" />
              Account Settings
            </Link>
            <Link
              to={'/settings' as never}
              onClick={() => setOpen(false)}
              className="flex items-center gap-2.5 rounded-lg px-2 py-1.5 text-[13px] transition-colors hover:bg-accent/50"
            >
              <BarChart3 className="h-4 w-4 text-muted-foreground" />
              Project Usage
            </Link>
          </div>

          <div className="my-1 h-px bg-border/50" />

          <div className="space-y-0.5 px-1 py-0.5">
            <button
              type="button"
              onClick={() => {
                setTheme(theme === 'dark' ? 'light' : 'dark');
                setOpen(false);
              }}
              className="flex w-full items-center gap-2.5 rounded-lg px-2 py-1.5 text-[13px] transition-colors hover:bg-accent/50"
            >
              {theme === 'dark' ? (
                <Sun className="h-4 w-4 text-muted-foreground" />
              ) : (
                <Moon className="h-4 w-4 text-muted-foreground" />
              )}
              {theme === 'dark' ? 'Light Theme' : 'Dark Theme'}
            </button>
          </div>

          <div className="my-1 h-px bg-border/50" />

          <div className="space-y-0.5 px-1 py-0.5">
            <button
              type="button"
              onClick={() => logout()}
              className="flex w-full items-center gap-2.5 rounded-lg px-2 py-1.5 text-[13px] text-destructive transition-colors hover:bg-destructive/10"
            >
              <LogOut className="h-4 w-4" />
              Log out
            </button>
          </div>
        </div>
      )}

      <button
        type="button"
        onClick={() => setOpen(!open)}
        className={`flex w-full items-center rounded-xl transition-colors hover:bg-sidebar-accent/50 ${
          collapsed ? 'justify-center gap-0 px-0 py-2' : 'gap-3 px-2.5 py-2'
        }`}
      >
        <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-primary/15 font-bold text-[10px] text-primary">
          {initials}
        </div>
        {!collapsed && (
          <>
            <div className="min-w-0 flex-1 text-left">
              <p className="truncate font-medium text-[12px] text-sidebar-foreground leading-tight">
                {user?.name ?? 'User'}
              </p>
            </div>
            <MoreVertical className="h-3.5 w-3.5 shrink-0 text-sidebar-foreground/40" />
          </>
        )}
      </button>
    </div>
  );
}
