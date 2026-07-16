import { cn } from '#/lib/utils';

interface EnvironmentSwitcherProps {
  value?: string;
  environments?: string[];
}

export function EnvironmentSwitcher({
  value = 'Production',
  environments = ['Production', 'Preview', 'Development'],
}: EnvironmentSwitcherProps) {
  return (
    <div className="inline-flex items-center rounded-md border bg-card p-1">
      {environments.map((environment) => {
        const selected = environment === value;

        return (
          <button
            key={environment}
            type="button"
            className={cn(
              'h-8 rounded px-3 font-medium text-xs transition-[background-color,color,box-shadow]',
              selected
                ? 'bg-primary text-primary-foreground shadow-sm'
                : 'text-muted-foreground hover:bg-muted hover:text-foreground'
            )}
          >
            {environment}
          </button>
        );
      })}
    </div>
  );
}
