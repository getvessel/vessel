import { Key01Icon, Logout01Icon, Settings01Icon, UserIcon } from '@hugeicons/core-free-icons';
import { HugeiconsIcon } from '@hugeicons/react';
import { createFileRoute } from '@tanstack/react-router';
import { useStore } from '@tanstack/react-store';
import { useLogout } from '#/hooks/useAuth';
import { authStore } from '#/stores/authStore';

export const Route = createFileRoute('/_workspace/settings')({
  component: SettingsPage,
});

const sections = [
  { id: 'profile', label: 'Profile', icon: UserIcon },
  { id: 'security', label: 'Security', icon: Key01Icon },
  { id: 'general', label: 'General', icon: Settings01Icon },
];

function SettingsPage() {
  const authState = useStore(authStore);
  const { mutateAsync: logout } = useLogout();
  const user = authState.user;

  const initials = user?.name
    ? user.name
        .split(' ')
        .map((n) => n[0])
        .join('')
        .toUpperCase()
    : 'U';

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h2 className="text-2xl font-bold tracking-tight text-foreground">Settings</h2>
        <p className="text-sm text-muted-foreground mt-1">
          Manage your account and workspace preferences.
        </p>
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-4">
        {/* Sidebar nav */}
        <nav className="flex flex-col gap-1">
          {sections.map((s) => (
            <button
              key={s.id}
              type="button"
              className="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium text-muted-foreground hover:text-foreground hover:bg-muted/50 transition-colors text-left"
            >
              <HugeiconsIcon icon={s.icon} className="h-4 w-4" />
              {s.label}
            </button>
          ))}
        </nav>

        {/* Content */}
        <div className="lg:col-span-3 space-y-4">
          {/* Profile card */}
          <div className="rounded-xl border border-border bg-card p-6">
            <h3 className="text-sm font-semibold text-foreground mb-5">Profile</h3>
            <div className="flex items-center gap-4 mb-6">
              <div className="flex h-14 w-14 items-center justify-center rounded-2xl bg-primary/20 text-lg font-bold text-primary">
                {initials}
              </div>
              <div>
                <p className="text-sm font-semibold text-foreground">{user?.name ?? 'User'}</p>
                <p className="text-xs text-muted-foreground">{user?.email ?? ''}</p>
              </div>
              <button
                type="button"
                className="ml-auto text-xs font-medium text-primary hover:underline"
              >
                Change avatar
              </button>
            </div>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div className="space-y-1.5">
                <label htmlFor="fullname" className="text-xs font-medium text-muted-foreground">Full name</label>
                <input
                  id="fullname"
                  type="text"
                  defaultValue={user?.name ?? ''}
                  className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary transition-all"
                />
              </div>
              <div className="space-y-1.5">
                <label htmlFor="email" className="text-xs font-medium text-muted-foreground">Email</label>
                <input
                  id="email"
                  type="email"
                  defaultValue={user?.email ?? ''}
                  className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary transition-all"
                />
              </div>
            </div>
            <div className="mt-4 flex justify-end">
              <button
                type="button"
                className="rounded-lg bg-primary px-4 py-2 text-xs font-semibold text-primary-foreground hover:bg-primary/90 transition-colors"
              >
                Save changes
              </button>
            </div>
          </div>

          {/* Danger zone */}
          <div className="rounded-xl border border-destructive/30 bg-destructive/5 p-6">
            <h3 className="text-sm font-semibold text-destructive mb-1">Danger Zone</h3>
            <p className="text-xs text-muted-foreground mb-4">
              These actions are irreversible. Please proceed with caution.
            </p>
            <button
              type="button"
              onClick={() => logout()}
              className="flex items-center gap-2 rounded-lg border border-destructive/40 bg-background px-4 py-2 text-xs font-semibold text-destructive hover:bg-destructive/10 transition-colors"
            >
              <HugeiconsIcon icon={Logout01Icon} className="h-3.5 w-3.5" />
              Log out of all sessions
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
