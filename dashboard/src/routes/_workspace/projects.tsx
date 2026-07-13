import {
  Activity01Icon,
  Folder01Icon,
  GridViewIcon,
  LeftToRightListBulletIcon,
  AddCircleIcon,
  Search01Icon,
} from '@hugeicons/core-free-icons';
import { HugeiconsIcon } from '@hugeicons/react';
import { createFileRoute } from '@tanstack/react-router';
import { useState } from 'react';

export const Route = createFileRoute('/_workspace/projects')({
  component: ProjectsPage,
});

function ProjectsPage() {
  const [view, setView] = useState<'grid' | 'list'>('grid');

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold tracking-tight text-foreground">Projects</h2>
          <p className="text-sm text-muted-foreground mt-1">
            Deploy and manage your applications and services.
          </p>
        </div>
        <button
          type="button"
          className="flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground hover:bg-primary/90 transition-colors"
        >
          <HugeiconsIcon icon={AddCircleIcon} className="h-4 w-4" />
          New Project
        </button>
      </div>

      {/* Toolbar */}
      <div className="flex items-center justify-between gap-4">
        <div className="flex items-center gap-2 rounded-lg border border-border bg-muted/30 px-3 py-2 flex-1 max-w-sm">
          <HugeiconsIcon
            icon={Search01Icon}
            className="h-4 w-4 text-muted-foreground shrink-0"
          />
          <input
            type="text"
            placeholder="Search projects..."
            className="bg-transparent text-sm text-foreground placeholder:text-muted-foreground outline-none w-full"
          />
        </div>
        <div className="flex items-center gap-1 rounded-lg border border-border bg-card p-1">
          <button
            type="button"
            onClick={() => setView('grid')}
            className={`flex h-7 w-7 items-center justify-center rounded-md transition-colors ${view === 'grid' ? 'bg-primary text-primary-foreground' : 'text-muted-foreground hover:text-foreground'}`}
          >
            <HugeiconsIcon icon={GridViewIcon} className="h-4 w-4" />
          </button>
          <button
            type="button"
            onClick={() => setView('list')}
            className={`flex h-7 w-7 items-center justify-center rounded-md transition-colors ${view === 'list' ? 'bg-primary text-primary-foreground' : 'text-muted-foreground hover:text-foreground'}`}
          >
            <HugeiconsIcon icon={LeftToRightListBulletIcon} className="h-4 w-4" />
          </button>
        </div>
      </div>

      {/* Empty state */}
      <div className="rounded-xl border border-border bg-card flex flex-col items-center justify-center py-24 px-6 text-center">
        <div className="flex h-14 w-14 items-center justify-center rounded-2xl bg-primary/10 mb-5">
          <HugeiconsIcon icon={Folder01Icon} className="h-7 w-7 text-primary" />
        </div>
        <h3 className="text-base font-semibold text-foreground mb-1">No projects yet</h3>
        <p className="text-sm text-muted-foreground max-w-xs mb-6">
          Create your first project to start deploying applications and managing services.
        </p>
        <button
          type="button"
          className="flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground hover:bg-primary/90 transition-colors"
        >
          <HugeiconsIcon icon={AddCircleIcon} className="h-4 w-4" />
          Create Project
        </button>
        <div className="mt-8 flex items-center gap-6 text-xs text-muted-foreground">
          <div className="flex items-center gap-1.5">
            <HugeiconsIcon icon={Activity01Icon} className="h-3.5 w-3.5" />
            Auto deployments
          </div>
          <div className="flex items-center gap-1.5">
            <HugeiconsIcon icon={Folder01Icon} className="h-3.5 w-3.5" />
            Monorepo support
          </div>
        </div>
      </div>
    </div>
  );
}
