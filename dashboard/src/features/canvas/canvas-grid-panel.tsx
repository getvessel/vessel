import { Box, Database, GitBranch, Plus, Route, Server } from 'lucide-react';
import { Button } from '#/components/ui/button';

const nodeTypes = [
  { label: 'Service', icon: Server },
  { label: 'Database', icon: Database },
  { label: 'Storage', icon: Box },
  { label: 'Route', icon: Route },
];

export function CanvasGridPanel() {
  return (
    <div className="dashboard-grid relative min-h-[460px] overflow-hidden rounded-lg border bg-card">
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_center,transparent_0%,var(--background)_82%)] opacity-70" />
      <div className="relative flex min-h-[460px] flex-col">
        <div className="flex flex-wrap items-center justify-between gap-3 border-b bg-background/80 px-4 py-3 backdrop-blur-sm">
          <div>
            <p className="font-medium text-sm">Environment map</p>
            <p className="text-muted-foreground text-xs">Drag resources here to define topology.</p>
          </div>
          <Button size="sm">
            <Plus className="size-4" />
            Add resource
          </Button>
        </div>

        <div className="grid flex-1 place-items-center p-8">
          <div className="max-w-lg rounded-lg border bg-background/90 p-5 text-center shadow-sm backdrop-blur-sm">
            <div className="mx-auto mb-4 flex size-10 items-center justify-center rounded-md border bg-muted/30">
              <GitBranch className="size-5 text-muted-foreground" />
            </div>
            <h2 className="font-semibold text-base">Canvas is ready</h2>
            <p className="mt-2 text-muted-foreground text-sm leading-6">
              Add services, databases, storage, or routes to build a visual map of this project
              environment.
            </p>
            <div className="mt-5 grid grid-cols-2 gap-2 sm:grid-cols-4">
              {nodeTypes.map((node) => (
                <button
                  key={node.label}
                  type="button"
                  className="flex flex-col items-center gap-2 rounded-md border bg-card px-3 py-3 text-muted-foreground text-xs transition-colors hover:border-primary/40 hover:text-foreground"
                >
                  <node.icon className="size-4" />
                  {node.label}
                </button>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
