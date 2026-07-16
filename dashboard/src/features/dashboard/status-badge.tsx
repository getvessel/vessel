import { AlertCircle, CheckCircle2, Clock3, MinusCircle } from 'lucide-react';
import { Badge } from '#/components/ui/badge';

type StatusTone = 'healthy' | 'info' | 'warning' | 'danger' | 'neutral';

interface StatusBadgeProps {
  label: string;
  tone?: StatusTone;
}

const toneClassNames: Record<StatusTone, string> = {
  healthy:
    'border-emerald-200 bg-emerald-50 text-emerald-700 dark:border-emerald-900 dark:bg-emerald-950 dark:text-emerald-300',
  info: 'border-primary/25 bg-primary/10 text-primary',
  warning:
    'border-yellow-200 bg-yellow-50 text-yellow-800 dark:border-yellow-900 dark:bg-yellow-950 dark:text-yellow-300',
  danger:
    'border-red-200 bg-red-50 text-red-700 dark:border-red-900 dark:bg-red-950 dark:text-red-300',
  neutral: 'border-border bg-muted text-muted-foreground',
};

export function StatusBadge({ label, tone = 'neutral' }: StatusBadgeProps) {
  const Icon =
    tone === 'healthy'
      ? CheckCircle2
      : tone === 'warning'
        ? Clock3
        : tone === 'danger'
          ? AlertCircle
          : MinusCircle;

  return (
    <Badge variant="outline" className={toneClassNames[tone]}>
      <Icon className="size-3" />
      {label}
    </Badge>
  );
}
