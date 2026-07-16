import { ExternalLink, GitBranch, GitCommit, RotateCcw } from 'lucide-react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { StatusBadge } from '#/features/dashboard/status-badge';
import type { Deployment } from '#/interfaces/deployment';

interface DeployPreviewCardProps {
  deployment: Deployment;
}

export function DeployPreviewCard({ deployment }: DeployPreviewCardProps) {
  const statusTone = deployment.status === 'FAILED' ? 'danger' : 'healthy';
  const commitLabel =
    deployment.commitMessage || deployment.commitHash?.slice(0, 7) || 'Manual deployment';
  const queueAction = (label: string) => {
    toast.info(`${label} queued`, {
      description: `Deployment ${deployment.id} will run this action when connected.`,
    });
  };

  return (
    <article className="rounded-lg border bg-card p-4">
      <div className="flex flex-col gap-3 lg:flex-row lg:items-start lg:justify-between">
        <div className="min-w-0">
          <div className="flex items-center gap-2">
            <StatusBadge label={deployment.status} tone={statusTone} />
            <span className="text-muted-foreground text-xs">
              {new Date(deployment.updatedAt).toLocaleString()}
            </span>
          </div>
          <h2 className="mt-3 truncate font-semibold text-base">{commitLabel}</h2>
          <div className="mt-2 flex flex-wrap gap-3 text-muted-foreground text-xs">
            <span className="inline-flex items-center gap-1.5">
              <GitBranch className="size-3.5" />
              {deployment.branch || 'main'}
            </span>
            <span className="inline-flex items-center gap-1.5">
              <GitCommit className="size-3.5" />
              {deployment.commitHash?.slice(0, 7) || deployment.trigger || 'dashboard'}
            </span>
          </div>
        </div>
        <div className="flex shrink-0 gap-2">
          <Button type="button" variant="outline" size="sm" onClick={() => queueAction('Preview')}>
            <ExternalLink className="size-4" />
            Preview
          </Button>
          <Button type="button" variant="outline" size="sm" onClick={() => queueAction('Rollback')}>
            <RotateCcw className="size-4" />
            Rollback
          </Button>
        </div>
      </div>
    </article>
  );
}
