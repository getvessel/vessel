import type { LucideIcon } from 'lucide-react';
import type { ReactNode } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { Skeleton } from '#/components/ui/skeleton';
import { cn } from '#/lib/utils';
import { StatusBadge } from './status-badge';

interface PageAction {
  label: string;
  icon?: LucideIcon;
  variant?: 'default' | 'outline' | 'secondary' | 'ghost' | 'destructive';
  disabled?: boolean;
  onClick?: () => void;
}

interface PageMetric {
  label: string;
  value?: string;
  detail: string;
  icon: LucideIcon;
  tone?: 'info' | 'healthy' | 'warning' | 'neutral';
}

interface OperationalPageProps {
  title: string;
  description: string;
  scope?: string;
  statusLabel?: string;
  statusTone?: 'healthy' | 'info' | 'warning' | 'danger' | 'neutral';
  primaryAction?: PageAction;
  secondaryActions?: PageAction[];
  metrics?: PageMetric[];
  children: ReactNode;
}

export function OperationalPage({
  title,
  description,
  scope,
  statusLabel = 'Ready',
  statusTone = 'neutral',
  primaryAction,
  secondaryActions = [],
  metrics = [],
  children,
}: OperationalPageProps) {
  const metricGridStyle = {
    gridTemplateColumns: 'repeat(auto-fit, minmax(min(100%, 260px), 1fr))',
  };

  const handleAction = (action: PageAction) => {
    if (action.onClick) {
      action.onClick();
      return;
    }

    toast.info(`${action.label} queued`, {
      description: 'This action will run as soon as the workflow endpoint is available.',
    });
  };

  return (
    <div className="flex flex-col gap-5">
      <div className="flex flex-col gap-4 border-b pb-4">
        <div className="flex flex-col gap-3 lg:flex-row lg:items-start lg:justify-between">
          <div className="min-w-0 space-y-2">
            <div className="flex flex-wrap items-center gap-2">
              {scope ? (
                <span className="font-medium text-muted-foreground text-xs uppercase tracking-wide">
                  {scope}
                </span>
              ) : null}
              <StatusBadge label={statusLabel} tone={statusTone} />
            </div>
            <div className="space-y-1">
              <h1 className="font-semibold text-xl tracking-tight">{title}</h1>
              <p className="max-w-3xl text-muted-foreground text-sm leading-6">{description}</p>
            </div>
          </div>
          <div className="flex shrink-0 flex-wrap gap-2">
            {secondaryActions.map((action) => (
              <Button
                key={action.label}
                type="button"
                variant={action.variant ?? 'outline'}
                disabled={action.disabled}
                onClick={() => handleAction(action)}
              >
                {action.icon ? <action.icon className="size-4" /> : null}
                {action.label}
              </Button>
            ))}
            {primaryAction ? (
              <Button
                type="button"
                variant={primaryAction.variant ?? 'default'}
                disabled={primaryAction.disabled}
                onClick={() => handleAction(primaryAction)}
              >
                {primaryAction.icon ? <primaryAction.icon className="size-4" /> : null}
                {primaryAction.label}
              </Button>
            ) : null}
          </div>
        </div>

        {metrics.length > 0 ? (
          <div className="grid gap-3" style={metricGridStyle}>
            {metrics.map((metric) => (
              <div
                key={metric.label}
                className="group relative overflow-hidden rounded-lg border bg-card p-4 transition-[background-color,border-color,transform] duration-150 hover:-translate-y-0.5 hover:border-sidebar-border hover:bg-card/95"
              >
                <div className="flex items-start justify-between gap-4">
                  <div
                    className={cn(
                      'flex size-9 shrink-0 items-center justify-center rounded-md border',
                      metric.tone === 'healthy' && 'border-success/20 bg-success/10 text-success',
                      metric.tone === 'warning' && 'border-warning/25 bg-warning/10 text-warning',
                      metric.tone === 'info' && 'border-primary/25 bg-primary/10 text-primary',
                      (!metric.tone || metric.tone === 'neutral') &&
                        'border-sidebar-border bg-muted/30 text-muted-foreground'
                    )}
                  >
                    <metric.icon className="size-4" />
                  </div>
                  <span
                    className={cn(
                      'mt-0.5 h-2 w-2 rounded-full',
                      metric.tone === 'healthy' && 'bg-success',
                      metric.tone === 'warning' && 'bg-warning',
                      metric.tone === 'info' && 'bg-primary',
                      (!metric.tone || metric.tone === 'neutral') && 'bg-muted-foreground/35'
                    )}
                  />
                </div>
                <div className="mt-4 min-w-0 space-y-1">
                  <p className="font-medium text-muted-foreground text-xs">{metric.label}</p>
                  {metric.value ? (
                    <p className="truncate font-semibold text-lg tracking-tight">{metric.value}</p>
                  ) : (
                    <Skeleton className="mt-1 h-5 w-16" />
                  )}
                  <p className="truncate text-muted-foreground text-xs leading-5">
                    {metric.detail}
                  </p>
                </div>
              </div>
            ))}
          </div>
        ) : null}
      </div>

      {children}
    </div>
  );
}
