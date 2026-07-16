import type { LucideIcon } from 'lucide-react';
import { Progress } from '#/components/ui/progress';

export function PressureRow({
  icon: Icon,
  label,
  value,
  detail,
}: {
  icon: LucideIcon;
  label: string;
  value?: number;
  detail?: string;
}) {
  const boundedValue = Math.min(Math.max(value ?? 0, 0), 100);

  return (
    <div className="grid gap-2 rounded-md border bg-muted/10 p-3">
      <div className="flex items-center justify-between gap-3">
        <div className="flex items-center gap-2">
          <Icon className="size-4 text-muted-foreground" />
          <span className="font-medium text-sm">{label}</span>
        </div>
        <span className="font-medium text-sm">
          {value === undefined ? 'Loading' : `${Math.round(value)}%`}
        </span>
      </div>
      <Progress value={boundedValue} />
      <p className="text-muted-foreground text-xs">{detail ?? 'Waiting for daemon metrics'}</p>
    </div>
  );
}
