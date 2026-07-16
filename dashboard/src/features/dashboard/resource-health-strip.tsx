import type { LucideIcon } from 'lucide-react';
import { cn } from '#/lib/utils';

interface HealthItem {
  label: string;
  value: string;
  detail: string;
  icon: LucideIcon;
  tone?: 'healthy' | 'warning' | 'danger' | 'info' | 'neutral';
}

interface ResourceHealthStripProps {
  items: HealthItem[];
}

const toneClasses = {
  healthy: 'bg-success',
  warning: 'bg-warning',
  danger: 'bg-destructive',
  info: 'bg-primary',
  neutral: 'bg-muted-foreground/40',
};

export function ResourceHealthStrip({ items }: ResourceHealthStripProps) {
  return (
    <div className="grid gap-2 rounded-lg border bg-card p-2 md:grid-cols-2 xl:grid-cols-4">
      {items.map((item) => (
        <div key={item.label} className="flex min-w-0 items-center gap-3 rounded-md px-3 py-2">
          <div className="flex size-8 shrink-0 items-center justify-center rounded-md border bg-muted/25">
            <item.icon className="size-4 text-muted-foreground" />
          </div>
          <div className="min-w-0 flex-1">
            <div className="flex items-center gap-2">
              <span className={cn('size-1.5 rounded-full', toneClasses[item.tone ?? 'neutral'])} />
              <p className="truncate font-medium text-muted-foreground text-xs">{item.label}</p>
            </div>
            <p className="truncate font-semibold text-sm">{item.value}</p>
            <p className="truncate text-muted-foreground text-xs">{item.detail}</p>
          </div>
        </div>
      ))}
    </div>
  );
}
