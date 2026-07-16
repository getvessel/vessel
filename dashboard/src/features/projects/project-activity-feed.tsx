import { GitCommit, Globe2, Rocket, Settings } from 'lucide-react';

const activityItems = [
  {
    label: 'Deployment promoted',
    detail: 'main deployed to production',
    time: '4m ago',
    icon: Rocket,
  },
  {
    label: 'Domain checked',
    detail: 'TLS certificate remains valid',
    time: '18m ago',
    icon: Globe2,
  },
  {
    label: 'Environment updated',
    detail: 'Preview variables synchronized',
    time: '42m ago',
    icon: Settings,
  },
  {
    label: 'Commit received',
    detail: 'feat: improve project importer',
    time: '1h ago',
    icon: GitCommit,
  },
];

export function ProjectActivityFeed() {
  return (
    <section className="rounded-lg border bg-card">
      <div className="border-b p-4">
        <h2 className="font-semibold text-base">Activity</h2>
        <p className="text-muted-foreground text-sm">Recent changes across this workspace.</p>
      </div>
      <div className="divide-y">
        {activityItems.map((item) => (
          <div key={item.label} className="flex items-start gap-3 p-4">
            <div className="flex size-8 shrink-0 items-center justify-center rounded-md border bg-muted/25">
              <item.icon className="size-4 text-muted-foreground" />
            </div>
            <div className="min-w-0 flex-1">
              <p className="font-medium text-sm">{item.label}</p>
              <p className="truncate text-muted-foreground text-sm">{item.detail}</p>
            </div>
            <span className="shrink-0 text-muted-foreground text-xs">{item.time}</span>
          </div>
        ))}
      </div>
    </section>
  );
}
