import { createFileRoute, Link } from '@tanstack/react-router';
import {
  Activity,
  ArrowRight,
  Boxes,
  Database,
  FolderKanban,
  HardDrive,
  MemoryStick,
  RefreshCw,
} from 'lucide-react';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '#/components/ui/card';
import { formatDuration, formatMb } from '#/features/dashboard/dashboard-format';
import { MetricCard } from '#/features/dashboard/metric-card';
import { StatusLine } from '#/features/dashboard/status-line';
import { SystemPressurePanel } from '#/features/dashboard/system-pressure-panel';
import { useListDatabases } from '#/hooks/useDatabases';
import { useListProjects } from '#/hooks/useProjects';
import { useGetSystemStats } from '#/hooks/useSystem';

export const Route = createFileRoute('/_dashboard/')({
  component: DashboardPage,
});

const quickLinks = [
  {
    title: 'Inspect projects',
    description: 'Apps, environments, domains, and settings.',
    to: '/projects',
    icon: FolderKanban,
  },
  {
    title: 'Review databases',
    description: 'Inventory, backups, SQL, and data tools.',
    to: '/databases',
    icon: Database,
  },
  {
    title: 'Import workloads',
    description: 'Bring Railway or Vercel projects in.',
    to: '/imports/railway',
    icon: Boxes,
  },
  {
    title: 'Tune settings',
    description: 'DNS, updates, users, Git apps, and cleanup.',
    to: '/settings',
    icon: HardDrive,
  },
];

function DashboardPage() {
  const systemStats = useGetSystemStats();
  const projects = useListProjects();
  const databases = useListDatabases();

  const stats = systemStats.data?.data;
  const projectCount = projects.data?.data.total ?? projects.data?.data.records.length ?? 0;
  const databaseCount = databases.data?.data.length ?? 0;
  const runningDatabases =
    databases.data?.data.filter((database) => database.status === 'running').length ?? 0;
  const isLoading = systemStats.isLoading || projects.isLoading || databases.isLoading;
  const hasError = systemStats.isError || projects.isError || databases.isError;
  const pressureState =
    stats && Math.max(stats.cpu.percent, stats.memory.percent, stats.disk.percent) > 85
      ? 'High pressure'
      : 'Nominal';

  return (
    <div className="flex flex-col gap-5">
      <div className="flex flex-col gap-3 border-b pb-5 lg:flex-row lg:items-center lg:justify-between">
        <Badge variant={hasError ? 'destructive' : 'secondary'} className="w-fit">
          {hasError ? 'Needs attention' : 'Live'}
        </Badge>
        <div className="flex flex-wrap gap-2">
          <Button
            variant="outline"
            onClick={() => {
              systemStats.refetch();
              projects.refetch();
              databases.refetch();
            }}
          >
            <RefreshCw className="size-4" />
            Refresh
          </Button>
          <Button asChild>
            <Link to="/projects">
              <FolderKanban className="size-4" />
              Open projects
            </Link>
          </Button>
        </div>
      </div>

      <div className="grid gap-3 md:grid-cols-2 xl:grid-cols-4">
        <MetricCard
          icon={FolderKanban}
          label="Projects"
          value={isLoading ? undefined : projectCount.toString()}
          detail="Tracked workspaces"
        />
        <MetricCard
          icon={Database}
          label="Databases"
          value={isLoading ? undefined : databaseCount.toString()}
          detail={`${runningDatabases} running`}
        />
        <MetricCard
          icon={Activity}
          label="CPU"
          value={stats ? `${Math.round(stats.cpu.percent)}%` : undefined}
          detail={stats ? `${stats.cpu.cores} cores` : 'Host utilization'}
        />
        <MetricCard
          icon={MemoryStick}
          label="Memory"
          value={stats ? `${Math.round(stats.memory.percent)}%` : undefined}
          detail={
            stats
              ? `${formatMb(stats.memory.usedMb)} / ${formatMb(stats.memory.totalMb)}`
              : 'RAM usage'
          }
        />
      </div>

      <div className="grid items-start gap-4 xl:grid-cols-[minmax(0,1fr)_420px]">
        <SystemPressurePanel stats={stats} state={pressureState} />

        <Card className="rounded-lg shadow-none">
          <CardHeader className="border-b">
            <CardTitle className="text-base">Daemon status</CardTitle>
            <CardDescription>Backend API and host process signals.</CardDescription>
            <CardAction>
              <Activity className="size-4 text-muted-foreground" />
            </CardAction>
          </CardHeader>
          <CardContent className="grid gap-4 text-sm">
            <StatusLine label="API" value={systemStats.isError ? 'Unavailable' : 'Reachable'} />
            <StatusLine
              label="Processes"
              value={stats ? stats.processes.toLocaleString() : 'Loading'}
            />
            <StatusLine
              label="Uptime"
              value={stats ? formatDuration(stats.uptimeSeconds) : 'Loading'}
            />
            <StatusLine
              label="Load average"
              value={stats ? stats.loadAvg.map((load) => load.toFixed(2)).join(' / ') : 'Loading'}
            />
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-3 md:grid-cols-2 xl:grid-cols-4">
        {quickLinks.map((link) => (
          <Link
            key={link.to}
            to={link.to as never}
            className="group flex min-h-24 items-center justify-between gap-4 rounded-lg border bg-card p-4 transition-colors hover:border-primary/40 hover:bg-muted/20"
          >
            <div className="flex min-w-0 items-center gap-3">
              <div className="flex size-9 items-center justify-center rounded-md border bg-muted/20">
                <link.icon className="size-4 text-muted-foreground group-hover:text-primary" />
              </div>
              <div className="min-w-0">
                <p className="font-medium text-sm">{link.title}</p>
                <p className="mt-1 line-clamp-2 text-muted-foreground text-xs">
                  {link.description}
                </p>
              </div>
            </div>
            <ArrowRight className="size-4 shrink-0 text-muted-foreground transition-transform group-hover:translate-x-0.5 group-hover:text-primary" />
          </Link>
        ))}
      </div>
    </div>
  );
}
