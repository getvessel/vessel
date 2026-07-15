import { Link } from '@tanstack/react-router';
import {
  BarChart3,
  FileText,
  HelpCircle,
  LogOut,
  Moon,
  MoreVertical,
  Sun,
  UserCircle,
} from 'lucide-react';
import { useTheme } from 'next-themes';
import { useEffect, useRef, useState } from 'react';
import { useLogout } from '#/hooks/useAuth';
import { useAuthState } from '#/stores/authStore';

export function UserMenu() {
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
    <div ref={ref} className="relative border-sidebar-border border-t p-3">
      {open && (
        <div className="fade-in zoom-in-95 absolute right-3 bottom-full left-3 z-50 mb-2 animate-in overflow-hidden rounded-lg border border-border bg-popover py-1 shadow-xl duration-100">
          <div className="mb-1 border-border border-b px-3 py-2">
            <p className="font-semibold text-[13px] text-foreground">{user?.name}</p>
            <p className="text-[11px] text-muted-foreground">{user?.email}</p>
          </div>

          <div className="space-y-0.5 px-1 py-1">
            <Link
              to={'/settings' as never}
              onClick={() => setOpen(false)}
              className="flex items-center gap-2.5 rounded-md px-2 py-1.5 text-[13px] transition-colors hover:bg-accent"
            >
              <UserCircle className="h-4 w-4 text-muted-foreground" />
              Account Settings
            </Link>
            <Link
              to={'/settings' as never}
              onClick={() => setOpen(false)}
              className="flex items-center gap-2.5 rounded-md px-2 py-1.5 text-[13px] transition-colors hover:bg-accent"
            >
              <BarChart3 className="h-4 w-4 text-muted-foreground" />
              Project Usage
            </Link>
          </div>

          <div className="my-1 h-px bg-border" />

          <div className="space-y-0.5 px-1 py-1">
            <Link
              to={'/docs' as never}
              onClick={() => setOpen(false)}
              className="flex items-center gap-2.5 rounded-md px-2 py-1.5 text-[13px] transition-colors hover:bg-accent"
            >
              <FileText className="h-4 w-4 text-muted-foreground" />
              Documentation
            </Link>
            <Link
              to={'/support' as never}
              onClick={() => setOpen(false)}
              className="flex items-center gap-2.5 rounded-md px-2 py-1.5 text-[13px] transition-colors hover:bg-accent"
            >
              <HelpCircle className="h-4 w-4 text-muted-foreground" />
              Support
            </Link>
          </div>

          <div className="my-1 h-px bg-border" />

          <div className="space-y-0.5 px-1 py-1">
            <button
              type="button"
              onClick={() => {
                setTheme(theme === 'dark' ? 'light' : 'dark');
                setOpen(false);
              }}
              className="flex w-full items-center gap-2.5 rounded-md px-2 py-1.5 text-[13px] transition-colors hover:bg-accent"
            >
              {theme === 'dark' ? (
                <Sun className="h-4 w-4 text-muted-foreground" />
              ) : (
                <Moon className="h-4 w-4 text-muted-foreground" />
              )}
              {theme === 'dark' ? 'Light Theme' : 'Dark Theme'}
            </button>
            <button
              type="button"
              onClick={() => logout()}
              className="flex w-full items-center gap-2.5 rounded-md px-2 py-1.5 text-[13px] text-destructive transition-colors hover:bg-destructive/10"
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
        className="flex w-full items-center gap-2.5 rounded-lg p-1.5 transition-colors hover:bg-sidebar-accent/50"
      >
        <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-primary/20 font-bold text-[11px] text-primary">
          {initials}
        </div>
        <div className="min-w-0 flex-1 text-left">
          <p className="truncate font-semibold text-[12px] text-sidebar-foreground leading-none">
            {user?.name ?? 'User'}
          </p>
        </div>
        <MoreVertical className="h-4 w-4 shrink-0 text-sidebar-foreground/50" />
      </button>
    </div>
  );
}
