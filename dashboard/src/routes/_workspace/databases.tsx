import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_workspace/databases')({
  component: DatabasesPage,
});

function DatabasesPage() {
  return (
    <div className="p-6">
      <h1 className="text-2xl font-semibold mb-4">Databases</h1>
      <p className="text-muted-foreground">Databases content goes here.</p>
    </div>
  );
}
