import { Database01Icon, AddCircleIcon } from '@hugeicons/core-free-icons';
import { HugeiconsIcon } from '@hugeicons/react';
import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_workspace/databases')({
  component: DatabasesPage,
});

const dbTypes = [
  {
    name: 'PostgreSQL',
    description: 'Relational database',
    color: 'text-blue-500',
    bg: 'bg-blue-500/10',
  },
  {
    name: 'MySQL',
    description: 'Relational database',
    color: 'text-orange-500',
    bg: 'bg-orange-500/10',
  },
  { name: 'Redis', description: 'In-memory cache', color: 'text-red-500', bg: 'bg-red-500/10' },
  {
    name: 'MongoDB',
    description: 'NoSQL document store',
    color: 'text-green-500',
    bg: 'bg-green-500/10',
  },
];

function DatabasesPage() {
  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold tracking-tight text-foreground">Databases & Storage</h2>
          <p className="text-sm text-muted-foreground mt-1">
            Manage your databases, caches, and object storage.
          </p>
        </div>
        <button
          type="button"
          className="flex items-center gap-2 rounded-lg bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground hover:bg-primary/90 transition-colors"
        >
          <HugeiconsIcon icon={AddCircleIcon} className="h-4 w-4" />
          New Database
        </button>
      </div>

      {/* Quick add section */}
      <div>
        <h3 className="mb-3 text-sm font-semibold text-foreground">Add a database</h3>
        <div className="grid grid-cols-2 gap-3 sm:grid-cols-4">
          {dbTypes.map((db) => (
            <button
              key={db.name}
              type="button"
              className="flex flex-col items-center gap-2 rounded-xl border border-border bg-card p-5 hover:border-primary/40 hover:shadow-sm transition-all duration-150 text-center"
            >
              <div className={`flex h-10 w-10 items-center justify-center rounded-xl ${db.bg}`}>
                <HugeiconsIcon icon={Database01Icon} className={`h-5 w-5 ${db.color}`} />
              </div>
              <div>
                <p className="text-sm font-semibold text-foreground">{db.name}</p>
                <p className="text-xs text-muted-foreground">{db.description}</p>
              </div>
            </button>
          ))}
        </div>
      </div>

      {/* Empty state for existing databases */}
      <div>
        <h3 className="mb-3 text-sm font-semibold text-foreground">Your databases</h3>
        <div className="rounded-xl border border-border bg-card flex flex-col items-center justify-center py-20 px-6 text-center">
          <div className="flex h-14 w-14 items-center justify-center rounded-2xl bg-violet-500/10 mb-5">
            <HugeiconsIcon icon={Database01Icon} className="h-7 w-7 text-violet-500" />
          </div>
          <h3 className="text-base font-semibold text-foreground mb-1">No databases yet</h3>
          <p className="text-sm text-muted-foreground max-w-xs">
            Choose a database type above to create your first managed database instance.
          </p>
        </div>
      </div>
    </div>
  );
}
