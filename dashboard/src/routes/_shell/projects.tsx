import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_shell/projects')({
  component: ProjectsPage,
});

function ProjectsPage() {
  return (
    <div className="p-6">
      <h1 className="text-2xl font-semibold mb-4">Projects</h1>
      <p className="text-muted-foreground">Projects content goes here.</p>
    </div>
  );
}
