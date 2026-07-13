import {
  Activity01Icon,
  CloudServerIcon,
  DashboardSquare01Icon,
  Database01Icon,
  Folder01Icon,
  HelpCircleIcon,
  Logout01Icon,
  Moon02Icon,
  Settings01Icon,
  Sun01Icon,
  UserGroupIcon,
} from "@hugeicons/core-free-icons";
import { HugeiconsIcon } from "@hugeicons/react";
import { Link, useRouterState } from "@tanstack/react-router";
import { useStore } from "@tanstack/react-store";
import { useTheme } from "next-themes";
import { useLogout } from "#/hooks/useAuth";
import { authStore } from "#/stores/authStore";

const navItems = [
  { title: "Dashboard", url: "/", icon: DashboardSquare01Icon, exact: true },
  { title: "Projects", url: "/projects", icon: Folder01Icon },
  { title: "Databases", url: "/databases", icon: Database01Icon },
  { title: "Deployments", url: "/deployments", icon: Activity01Icon },
  { title: "Teams", url: "/teams", icon: UserGroupIcon },
];

const bottomItems = [
  { title: "Settings", url: "/settings", icon: Settings01Icon },
  { title: "Support", url: "/support", icon: HelpCircleIcon },
];

function NavItem({
  item,
  exact = false,
}: {
  item: (typeof navItems)[0];
  exact?: boolean;
}) {
  const routerState = useRouterState();
  const pathname = routerState.location.pathname;
  const isActive = exact
    ? pathname === item.url
    : pathname.startsWith(item.url);

  return (
    <Link
      to={item.url}
      className={[
        "flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-all duration-150",
        isActive
          ? "bg-sidebar-accent text-sidebar-accent-foreground"
          : "text-sidebar-foreground/60 hover:text-sidebar-foreground hover:bg-sidebar-accent/50",
      ].join(" ")}
    >
      <HugeiconsIcon
        icon={item.icon}
        className={["h-4 w-4 shrink-0", isActive ? "text-primary" : ""].join(
          " ",
        )}
      />
      <span>{item.title}</span>
    </Link>
  );
}

export function AppSidebar() {
  const { theme, setTheme } = useTheme();
  const authState = useStore(authStore);
  const { mutateAsync: logout } = useLogout();
  const user = authState.user;

  const initials = user?.name
    ? user.name
        .split(" ")
        .map((n) => n[0])
        .join("")
        .toUpperCase()
    : "U";

  return (
    <aside className="fixed inset-y-0 left-0 z-20 flex w-55 flex-col border-r border-sidebar-border bg-sidebar">
      {/* Logo */}
      <div className="flex h-14 shrink-0 items-center gap-2.5 px-4 border-b border-sidebar-border">
        <div className="flex h-7 w-7 items-center justify-center rounded-lg bg-primary/10">
          <HugeiconsIcon
            icon={CloudServerIcon}
            className="h-4 w-4 text-primary"
          />
        </div>
        <span className="font-semibold text-[15px] tracking-tight text-sidebar-foreground">
          Vessl
        </span>
        <span className="ml-auto text-[10px] font-medium text-muted-foreground bg-muted/60 px-1.5 py-0.5 rounded">
          v0.1
        </span>
      </div>

      {/* Nav */}
      <nav className="flex flex-1 flex-col gap-0.5 overflow-y-auto p-3">
        <p className="mb-1 px-3 text-[10px] font-semibold uppercase tracking-widest text-sidebar-foreground/30">
          Platform
        </p>
        {navItems.map((item) => (
          <NavItem key={item.url} item={item} exact={item.exact} />
        ))}

        <div className="mt-4 mb-1 h-px bg-sidebar-border" />
        <p className="mb-1 px-3 text-[10px] font-semibold uppercase tracking-widest text-sidebar-foreground/30">
          Account
        </p>
        {bottomItems.map((item) => (
          <NavItem key={item.url} item={item} />
        ))}
      </nav>

      {/* Footer */}
      <div className="shrink-0 border-t border-sidebar-border p-3 space-y-1">
        {/* Theme toggle */}
        <button
          type="button"
          onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
          className="flex w-full items-center gap-3 rounded-md px-3 py-2 text-sm font-medium text-sidebar-foreground/60 hover:text-sidebar-foreground hover:bg-sidebar-accent/50 transition-all duration-150"
        >
          <HugeiconsIcon
            icon={theme === "dark" ? Sun01Icon : Moon02Icon}
            className="h-4 w-4 shrink-0"
          />
          <span>{theme === "dark" ? "Light mode" : "Dark mode"}</span>
        </button>

        {/* Logout */}
        <button
          type="button"
          onClick={() => logout()}
          className="flex w-full items-center gap-3 rounded-md px-3 py-2 text-sm font-medium text-destructive/70 hover:text-destructive hover:bg-destructive/10 transition-all duration-150"
        >
          <HugeiconsIcon icon={Logout01Icon} className="h-4 w-4 shrink-0" />
          <span>Log out</span>
        </button>

        {/* User info */}
        <div className="flex items-center gap-2.5 mt-2 px-3 py-2 rounded-md bg-sidebar-accent/40">
          <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-primary/20 text-[10px] font-bold text-primary">
            {initials}
          </div>
          <div className="min-w-0">
            <p className="truncate text-xs font-semibold text-sidebar-foreground">
              {user?.name ?? "User"}
            </p>
            <p className="truncate text-[10px] text-sidebar-foreground/50">
              {user?.email ?? ""}
            </p>
          </div>
        </div>
      </div>
    </aside>
  );
}
