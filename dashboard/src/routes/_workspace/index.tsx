import {
  Activity01Icon,
  AddCircleIcon,
  CloudServerIcon,
  Database01Icon,
  Folder01Icon,
  TrendingUpDownIcon,
} from '@hugeicons/core-free-icons';
import { HugeiconsIcon } from '@hugeicons/react';
import { createFileRoute, Link } from '@tanstack/react-router';
import { useStore } from '@tanstack/react-store';
import { authStore } from '#/stores/authStore';

export const Route = createFileRoute('/_workspace/')({
  component: DashboardPage,
});

const stats = [
  {
    label: 'Active Projects',
    value: '0',
    icon: Folder01Icon,
    color: 'text-blue-500',
    bg: 'bg-blue-500/10',
    change: '+0 this month',
  },
  {
    label: 'Databases',
    value: '0',
    icon: Database01Icon,
    color: 'text-violet-500',
    bg: 'bg-violet-500/10',
    change: '0 running',
  },
  {
    label: 'Deployments',
    value: '0',
    icon: Activity01Icon,
    color: 'text-emerald-500',
    bg: 'bg-emerald-500/10',
    change: '0 successful',
  },
  {
    label: 'Uptime',
    value: '99.9%',
    icon: TrendingUpDownIcon,
    color: 'text-amber-500',
    bg: 'bg-amber-500/10',
    change: 'Last 30 days',
  },
];

const quickActions = [
  {
    title: 'New Project',
    description: 'Deploy a new application or service',
    icon: Folder01Icon,
    href: '/projects',
    color: 'text-blue-500',
    bg: 'bg-blue-500/10 hover:bg-blue-500/20',
  },
  {
    title: 'Add Database',
    description: 'Spin up PostgreSQL, MySQL or Redis',
    icon: Database01Icon,
    href: '/databases',
    color: 'text-violet-500',
    bg: 'bg-violet-500/10 hover:bg-violet-500/20',
  },
  {
    title: 'Connect Server',
    description: 'Add a new self-hosted server',
    icon: CloudServerIcon,
    href: '/settings',
    color: 'text-emerald-500',
    bg: 'bg-emerald-500/10 hover:bg-emerald-500/20',
  },
];

function DashboardPage() {
  const authState = useStore(authStore);
  const user = authState.user;
  const firstName = user?.name?.split(' ')[0] ?? 'there';

  return (
    <div className="space-y-8">
      {/* Welcome header */}
      <div className="flex items-start justify-between">
        <div>
          <h2 className="text-2xl font-bold text-foreground tracking-tight">
            Good morning, {firstName} 👋
          </h2>
          <p className="mt-1 text-sm text-muted-foreground">
            Here's what's happening with your infrastructure today.
          </p>
        </div>
        <Link
          to="/projects"
          className="flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground hover:bg-primary/90 transition-colors duration-150"
        >
          <HugeiconsIcon icon={AddCircleIcon} className="h-4 w-4" />
          New Project
        </Link>
      </div>

      {/* Stats grid */}
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => (
          <div
            key={stat.label}
            className="rounded-xl border border-border bg-card p-5 transition-shadow hover:shadow-sm"
          >
            <div className="flex items-center justify-between">
              <p className="text-sm font-medium text-muted-foreground">{stat.label}</p>
              <div className={`flex h-8 w-8 items-center justify-center rounded-lg ${stat.bg}`}>
                <HugeiconsIcon icon={stat.icon} className={`h-4 w-4 ${stat.color}`} />
              </div>
            </div>
            <p className="mt-3 text-3xl font-bold tracking-tight text-foreground">{stat.value}</p>
            <p className="mt-1 text-xs text-muted-foreground">{stat.change}</p>
          </div>
        ))}
      </div>

      {/* Main content: Quick actions + Recent activity */}
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        {/* Quick actions */}
        <div className="lg:col-span-1">
          <h3 className="mb-3 text-sm font-semibold text-foreground">Quick Actions</h3>
          <div className="space-y-2">
            {quickActions.map((action) => (
              <Link
                key={action.title}
                to={action.href}
                className={`flex items-center gap-4 rounded-xl border border-border bg-card p-4 transition-all duration-150 hover:border-border/80 hover:shadow-sm ${action.bg.split(' ')[1]}`}
              >
                <div
                  className={`flex h-9 w-9 shrink-0 items-center justify-center rounded-lg ${action.bg.split(' ')[0]}`}
                >
                  <HugeiconsIcon icon={action.icon} className={`h-4.5 w-4.5 ${action.color}`} />
                </div>
                <div className="min-w-0">
                  <p className="text-sm font-semibold text-foreground">{action.title}</p>
                  <p className="text-xs text-muted-foreground truncate">{action.description}</p>
                </div>
              </Link>
            ))}
          </div>
        </div>

        {/* Recent activity */}
        <div className="lg:col-span-2">
          <h3 className="mb-3 text-sm font-semibold text-foreground">Recent Activity</h3>
          <div className="rounded-xl border border-border bg-card overflow-hidden">
            {/* Empty state */}
            <div className="flex flex-col items-center justify-center py-16 px-6 text-center">
              <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-muted/60 mb-4">
                <HugeiconsIcon icon={Activity01Icon} className="h-5 w-5 text-muted-foreground" />
              </div>
              <p className="text-sm font-semibold text-foreground mb-1">No activity yet</p>
              <p className="text-xs text-muted-foreground max-w-xs">
                Once you create projects and run deployments, all your activity will appear here.
              </p>
              <Link
                to="/projects"
                className="mt-4 flex items-center gap-1.5 rounded-lg bg-primary/10 px-3 py-1.5 text-xs font-semibold text-primary hover:bg-primary/20 transition-colors"
              >
                <HugeiconsIcon icon={AddCircleIcon} className="h-3.5 w-3.5" />
                Create your first project
              </Link>
            </div>
          </div>
        </div>
      </div>

      {/* Infrastructure status */}
      <div>
        <h3 className="mb-3 text-sm font-semibold text-foreground">Infrastructure Status</h3>
        <div className="rounded-xl border border-border bg-card p-5">
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
            {[
              {
                label: 'Docker',
                status: 'Online',
                color: 'text-emerald-500',
                dot: 'bg-emerald-500',
              },
              {
                label: 'API Gateway',
                status: 'Online',
                color: 'text-emerald-500',
                dot: 'bg-emerald-500',
              },
              {
                label: 'Database Cluster',
                status: 'Online',
                color: 'text-emerald-500',
                dot: 'bg-emerald-500',
              },
            ].map((s) => (
              <div key={s.label} className="flex items-center gap-3">
                <span
                  className={`h-2 w-2 rounded-full ${s.dot} ring-4 ring-emerald-500/20 shrink-0`}
                />
                <div>
                  <p className="text-xs font-medium text-foreground">{s.label}</p>
                  <p className={`text-xs font-semibold ${s.color}`}>{s.status}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
