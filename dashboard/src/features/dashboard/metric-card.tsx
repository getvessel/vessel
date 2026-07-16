import type { LucideIcon } from 'lucide-react';
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '#/components/ui/card';
import { Skeleton } from '#/components/ui/skeleton';

export function MetricCard({
  icon: Icon,
  label,
  value,
  detail,
}: {
  icon: LucideIcon;
  label: string;
  value?: string;
  detail: string;
}) {
  return (
    <Card className="gap-0 rounded-lg py-0 shadow-none">
      <CardHeader className="p-4">
        <CardAction>
          <Icon className="size-4 text-muted-foreground" />
        </CardAction>
        <CardDescription>{label}</CardDescription>
        <CardTitle className="text-xl">{value ?? <Skeleton className="h-7 w-16" />}</CardTitle>
      </CardHeader>
      <CardContent className="px-4 pb-4 text-muted-foreground text-xs">{detail}</CardContent>
    </Card>
  );
}
