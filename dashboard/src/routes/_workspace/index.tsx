import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_workspace/')({
  component: DashboardPage,
});

function DashboardPage() {
  return (
    <div className="p-6">
      <h1 className="text-2xl font-semibold mb-4">Dashboard</h1>
      <p className="text-muted-foreground">Dashboard content goes here.</p>
    </div>
  );
}
