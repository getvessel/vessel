import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_workspace/settings')({
  component: SettingsPage,
});

function SettingsPage() {
  return (
    <div className="p-6">
      <h1 className="text-2xl font-semibold mb-4">Settings</h1>
      <p className="text-muted-foreground">Settings content goes here.</p>
    </div>
  );
}
