import { CheckCircle2, Circle, Loader2, XCircle } from 'lucide-react';
import type { DeploymentStatus } from '#/interfaces/deployment';
import { cn } from '#/lib/utils';

interface DeploymentTimelineProps {
  status?: DeploymentStatus;
}

const steps = ['Queued', 'Cloning', 'Installing', 'Building', 'Releasing', 'Live'];

function getCurrentStep(status?: DeploymentStatus) {
  if (status === 'FAILED') {
    return 3;
  }

  if (status === 'pending') {
    return 0;
  }

  if (status === 'CLONING') {
    return 1;
  }

  if (status === 'BUILDING' || status === 'PULLING') {
    return 3;
  }

  if (status === 'READY' || status === 'ACTIVE' || status === 'SUCCESS') {
    return 5;
  }

  return 4;
}

export function DeploymentTimeline({ status }: DeploymentTimelineProps) {
  const currentStep = getCurrentStep(status);
  const failed = status === 'FAILED';

  return (
    <section className="rounded-lg border bg-card">
      <div className="border-b p-4">
        <h2 className="font-semibold text-base">Deployment timeline</h2>
        <p className="text-muted-foreground text-sm">
          Build and release progress for latest deploy.
        </p>
      </div>
      <div className="p-4">
        <div className="grid gap-4">
          {steps.map((step, index) => {
            const complete =
              index < currentStep || (!failed && index === currentStep && index === 5);
            const active = index === currentStep && !complete;
            const stepFailed = failed && index === currentStep;

            return (
              <div key={step} className="grid grid-cols-[24px_1fr_auto] items-start gap-3">
                <div className="flex flex-col items-center">
                  {stepFailed ? (
                    <XCircle className="size-5 text-destructive" />
                  ) : active ? (
                    <Loader2 className="size-5 animate-spin text-primary" />
                  ) : complete ? (
                    <CheckCircle2 className="size-5 text-success" />
                  ) : (
                    <Circle className="size-5 text-muted-foreground/45" />
                  )}
                  {index < steps.length - 1 ? (
                    <div
                      className={cn('mt-2 h-8 w-px', complete ? 'bg-success/50' : 'bg-border')}
                    />
                  ) : null}
                </div>
                <div>
                  <p className="font-medium text-sm">{step}</p>
                  <p className="text-muted-foreground text-xs">
                    {stepFailed
                      ? 'Needs attention'
                      : active
                        ? 'In progress'
                        : complete
                          ? 'Complete'
                          : 'Waiting'}
                  </p>
                </div>
                <span className="text-muted-foreground text-xs">{complete ? '<1m' : '--'}</span>
              </div>
            );
          })}
        </div>
      </div>
    </section>
  );
}
