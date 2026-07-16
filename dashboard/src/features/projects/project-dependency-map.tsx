import { Boxes, Database, Globe2, HardDrive } from 'lucide-react';
import type { AppService } from '#/interfaces/deployment';

interface ProjectDependencyMapProps {
  services: AppService[];
}

export function ProjectDependencyMap({ services }: ProjectDependencyMapProps) {
  const primaryService = services[0]?.name ?? 'Service';
  const domain = services[0]?.domain || 'Domain';

  return (
    <section className="rounded-lg border bg-card">
      <div className="border-b p-4">
        <h2 className="font-semibold text-base">Dependency map</h2>
        <p className="text-muted-foreground text-sm">Runtime relationships for this environment.</p>
      </div>
      <div className="overflow-x-auto p-4">
        <div className="grid min-w-[620px] grid-cols-[1fr_56px_1fr_56px_1fr] items-center">
          <div className="rounded-lg border bg-muted/20 p-4">
            <div className="flex items-center gap-2">
              <Boxes className="size-4 text-primary" />
              <p className="font-medium text-sm">{primaryService}</p>
            </div>
            <p className="mt-1 text-muted-foreground text-xs">Application workload</p>
          </div>
          <div className="h-px bg-border" />
          <div className="grid gap-3">
            <div className="rounded-lg border bg-muted/20 p-4">
              <div className="flex items-center gap-2">
                <Database className="size-4 text-success" />
                <p className="font-medium text-sm">Database</p>
              </div>
              <p className="mt-1 text-muted-foreground text-xs">Persistent storage</p>
            </div>
            <div className="rounded-lg border bg-muted/20 p-4">
              <div className="flex items-center gap-2">
                <HardDrive className="size-4 text-warning" />
                <p className="font-medium text-sm">Object storage</p>
              </div>
              <p className="mt-1 text-muted-foreground text-xs">Uploads and backups</p>
            </div>
          </div>
          <div className="h-px bg-border" />
          <div className="rounded-lg border bg-muted/20 p-4">
            <div className="flex items-center gap-2">
              <Globe2 className="size-4 text-primary" />
              <p className="font-medium text-sm">{domain}</p>
            </div>
            <p className="mt-1 text-muted-foreground text-xs">Public edge route</p>
          </div>
        </div>
      </div>
    </section>
  );
}
