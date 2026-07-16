import type { LucideIcon } from 'lucide-react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';

interface EmptyPanelProps {
  icon: LucideIcon;
  title: string;
  description: string;
  actionLabel?: string;
}

export function EmptyPanel({ icon: Icon, title, description, actionLabel }: EmptyPanelProps) {
  return (
    <div className="flex min-h-48 flex-col items-center justify-center rounded-lg border border-dashed bg-muted/15 p-8 text-center">
      <div className="mb-4 flex size-9 items-center justify-center rounded-md border bg-background">
        <Icon className="size-5 text-muted-foreground" />
      </div>
      <h2 className="font-semibold text-base">{title}</h2>
      <p className="mt-2 max-w-md text-muted-foreground text-sm leading-6">{description}</p>
      {actionLabel ? (
        <Button
          className="mt-5"
          type="button"
          variant="outline"
          onClick={() =>
            toast.info(`${actionLabel} queued`, {
              description: 'This workflow will start when the action endpoint is available.',
            })
          }
        >
          {actionLabel}
        </Button>
      ) : null}
    </div>
  );
}
