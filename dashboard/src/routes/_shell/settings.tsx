import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_shell/settings')({
  component: SettingsPage,
});

function SettingsPage() {
  return (
    <div className="p-6">
      <h1 className="mb-4 font-semibold text-2xl">Settings</h1>
      <p className="text-muted-foreground">Settings content goes here.</p>
    </div>
  );
}
