import { createFileRoute, Navigate } from '@tanstack/react-router';
import { SetupForm } from '#/features/auth';
import { useGetSetupStatus } from '#/hooks/useSettings';

export const Route = createFileRoute('/_auth/setup')({
  component: SetupPage,
});

function SetupPage() {
  const { data: setupStatus, isLoading } = useGetSetupStatus();

  if (!isLoading && !setupStatus?.data?.setupRequired) {
    return <Navigate to="/login" replace />;
  }

  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="flex flex-col space-y-2 text-center mb-8">
        <h1 className="text-3xl font-semibold tracking-tight text-foreground">Welcome to Vessl</h1>
        <p className="text-sm text-muted-foreground">
          Set up the initial owner account to manage your cluster.
        </p>
      </div>

      <SetupForm />
    </div>
  );
}
